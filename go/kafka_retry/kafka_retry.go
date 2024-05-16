package kafka_retry

import (
    "context"
    "fmt"
    "math/rand"
    "time"

    "github.com/cenkalti/backoff/v4"
    "github.com/segmentio/kafka-go"
)

type ProcessRetryHandler interface {
    Process(context.Context, kafka.Message, int) error
}

type MessageWithIndex struct {
    Msg   kafka.Message
    Index int
}

type ConsumerWithRetryOptions struct {
    Handler    ProcessRetryHandler
    Reader     *kafka.Reader
    MaxRetries int
    ProcessRetryQueue chan MessageWithIndex
    ReadRetryQueue  chan MessageWithIndex
    Backoff    backoff.BackOff
}

type SimpleHandler struct{}

func (h *SimpleHandler) Process(ctx context.Context, m kafka.Message, index int) error {
    // Simulated processing error remains here
    if rand.Float32() < 0.5 { // 80% chance to succeed
        fmt.Printf("Message #%d processed successfully: %s\n", index, string(m.Value))
        return nil
    }
    fmt.Printf("Simulated processing error for message #%d: %s\n", index, string(m.Key))
    return fmt.Errorf("simulated processing error for message #%d: %s", index, string(m.Key))
}

func NewConsumerWithRetry(ctx context.Context, options *ConsumerWithRetryOptions) {
    rand.Seed(time.Now().UnixNano())

    // Goroutine for processing retries
    go func() {
        for mi := range options.ProcessRetryQueue {
            retries := 0
            for retries < options.MaxRetries {
                err := options.Handler.Process(ctx, mi.Msg, mi.Index)
                if err == nil {
                    fmt.Printf("Successfully processed message #%d after %d retry(ies).\n", mi.Index, retries)
                    break
                }
                retries++
                fmt.Printf("Failed to process message #%d, retrying. Attempt %d.\n", mi.Index, retries)
                time.Sleep(options.Backoff.NextBackOff())
            }
            if retries >= options.MaxRetries {
                fmt.Printf("Max retries exceeded for processing message #%d. Giving up.\n", mi.Index)
            }
        }
    }()

    // Goroutine for reading retries 
    go func() {
        for msgWithIndex := range options.ReadRetryQueue {
            readRetries := 0
            messageIndex := msgWithIndex.Index
            for readRetries < options.MaxRetries {
                fmt.Printf("Attempting to read message #%d. Retry attempt %d.\n", messageIndex, readRetries)
                msg, err := options.Reader.ReadMessage(ctx)
                
            
                if err != nil {
                    time.Sleep(options.Backoff.NextBackOff())
                    readRetries++
                    continue
                }

                fmt.Printf("Successfully read message #%d after %d retry(ies).\n", messageIndex, readRetries)
                
                // Process the message after successful read
                processErr := options.Handler.Process(ctx, msg, messageIndex)
                if processErr != nil {
                    fmt.Printf("Processing error for message #%d, sending to RetryQueue.\n", messageIndex)
                    options.ProcessRetryQueue <- MessageWithIndex{Msg: msg, Index: messageIndex}
                }
                break
            }

            if readRetries >= options.MaxRetries {
                fmt.Printf("Max retries exceeded for reading message #%d. Giving up.\n", messageIndex)
            }
        }
    }()
    messageCounter := 1
    for {
        select {
        case <-ctx.Done():
            return
        default:
            
        }
        msg, err := options.Reader.ReadMessage(ctx)
            fmt.Printf("Error reading message : %v\n", ctx)
            fmt.Printf("[Debug] Attempting to read message. Current counter: %d\n", messageCounter)
            if err != nil {
                fmt.Printf("[Debug] Error encountered reading message #%d, sending to ReadQueue.\n", messageCounter)
                options.ReadRetryQueue <- MessageWithIndex{Msg: msg, Index: messageCounter}
                messageCounter++
                continue
            }
    
            fmt.Printf("[Debug] Message #%d read successfully. Sending to Handler for processing.\n", messageCounter)
            processErr := options.Handler.Process(ctx, msg, messageCounter)
            if processErr != nil {
                fmt.Printf("[Debug] Processing error for message #%d, sending to RetryQueue.\n", messageCounter)
                options.ProcessRetryQueue <- MessageWithIndex{Msg: msg, Index: messageCounter}
            }
            messageCounter++
    }
    
}




//Şimdiki kullanacağımız yöntemde mesajalrın sıralaması önemli değil. Sıralama önemli olduğu zamanlarda multiple consumers kullanılmalı.
//Burda hem read hataları alırken topic okunmaya devam edilecek, hem de processing hataları olduğunda main topic okunmaya devam edilecek. (goroutines sayesinde).
//Yani şu şekilde: Read'te hata alındığında read'in retryQueue suna yollanacak. Hata almayana kadar denenecek.
//Queue dayken hata çözülünce process'e yollanacak. O şey process edilirken hata alırsa process queue'
//ya yollanacak. Hata almayana kadar denenecek. Her yeni mesajda yeni GoRoutine çalışmaz. Hata alınınca queue elementi olacak goRoutine çalışır.

//Her tekrar deneme arasında "backoff" kadar zaman geçecek. Ve tekrar denerken bu tekrar denemenin başarılı olacağı kesin değil(hem process, hem read için). Başarılı olana kadar 
//O mesajın tekrar process edilmesi ya da read edilmesi tekrarlanacak. (Sadece process fail edilmesi simüle edilmiştir.)

//Goroutineler queuelardakileri handle ederler.

//SOURCE: https://alexandrecastro.tech/post/kafka-retries-go/

//Analogy : The robot (Kafka consumer) continuously checks for messages (like watching a basket for balls).
//When a message arrives, the robot tries to process it (play a mini-game with the ball).
//If the robot fails to process it correctly (win the game), it tries again, up to a set number of times. This is the retry logic.
//Between each retry, the robot waits for a short time, which gets longer with more failures. This is the backoff strategy.
//If the robot can't process the message after several tries, it moves the message to a different place (the DLQ), for special handling later.

//KAFKA BIG STREAM (MULTIPLE CONSUMERS MULTIPLE PRODUCERS):
//With multiple producers and consumers in a Kafka system:

//1. **Multiple Producers**: Different producers can send messages to the same or different topics in Kafka. Imagine several people throwing balls into the same or different baskets.
//2. **Multiple Consumers**: Each consumer might be responsible for reading messages from a specific topic or partition of a topic. It's like having several robots, each watching over its own basket.
//3. **Partitioning and Load Balancing**: Kafka divides the data into partitions, which helps in distributing the load. Each consumer in a consumer group reads from a specific partition, ensuring efficient processing and no overlap in message consumption.
//4. **Concurrent Message Processing**: Multiple consumers can process messages concurrently, increasing throughput. Each robot handles balls from its basket independently.
//5. **Retry Logic for Each Consumer**: If any consumer fails to process a message, its own retry logic  (like the robot retrying the mini-game) with backoff is applied. Each robot deals with its failures independently.
//6. **DLQ Handling**: Similarly, each consumer manages its DLQ. If a robot can't process a ball, it moves it to its DLQ.
//In summary, with multiple producers and consumers, Kafka efficiently handles a higher volume of messages, distributing them across consumers for parallel processing. Each consumer independently manages retries and failures.




