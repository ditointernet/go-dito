package pubsub_test

// import (
// 	"context"

// 	"testing"

// 	"github.com/ditointernet/go-dito/errors"
// 	"github.com/ditointernet/go-dito/pubsub"
// 	"github.com/ditointernet/go-dito/pubsub/mocks"

// 	ps "cloud.google.com/go/pubsub"

// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// var (
// 	ctrl *gomock.Controller

// 	topicM  *mocks.MockTopicer
// 	resultM *mocks.MockResultier

// 	ctx context.Context
// )

// type MessageSchemaM struct {
// }

// func (ms MessageSchemaM) MarshalJSON() ([]byte, error) {
// 	b := []byte("ABCD")
// 	return b, nil
// }

// type MessageSchemaWithErrorM struct {
// }

// func (mswe MessageSchemaWithErrorM) MarshalJSON() ([]byte, error) {
// 	return nil, errors.New("Error to marshal message data")
// }

// func TestPubSub(t *testing.T) {
// 	RegisterFailHandler(Fail)

// 	BeforeEach(func() {
// 		ctrl = gomock.NewController(GinkgoT())

// 		topicM = mocks.NewMockTopicer(ctrl)
// 		resultM = mocks.NewMockResultier(ctrl)

// 		ctx = context.Background()
// 	})

// 	AfterEach(func() {
// 		ctrl.Finish()
// 	})

// 	RunSpecs(t, "PubSub Suite")
// }

// var _ = Describe("PubSubClient", func() {
// 	Context("Publish", func() {
// 		var (
// 			pubsubClient pubsub.PubSubClient[MessageSchemaM]
// 			publishInput pubsub.PublishInput[MessageSchemaM]

// 			pubsubClientWithError pubsub.PubSubClient[MessageSchemaWithErrorM]
// 			publishInputWithError pubsub.PublishInput[MessageSchemaWithErrorM]

// 			errs []error
// 		)

// 		BeforeEach(func() {
// 			pubsubClient = pubsub.MustNewPubSubClient[MessageSchemaM](topicM)

// 			publishInput = pubsub.PublishInput[MessageSchemaM]{
// 				Data: MessageSchemaM{},
// 				Attributes: map[string]string{
// 					"key": "value",
// 				},
// 			}

// 			pubsubClientWithError = pubsub.MustNewPubSubClient[MessageSchemaWithErrorM](topicM)

// 			publishInputWithError = pubsub.PublishInput[MessageSchemaWithErrorM]{
// 				Data: MessageSchemaWithErrorM{},
// 				Attributes: map[string]string{
// 					"key": "value",
// 				},
// 			}
// 		})

// 		Context("Error cases", func() {
// 			When("fails to marshal the data", func() {
// 				It("returns the error list", func() {
// 					errs = pubsubClientWithError.Publish(ctx, publishInputWithError)
// 					Expect(errs).To(Equal([]error{errors.New("Error to marshal message data")}))
// 				})
// 			})

// 			When("fails to publish message into pubsub topic", func() {
// 				It("returns the error list", func() {
// 					data, _ := publishInput.Data.MarshalJSON()

// 					pubSubMsg := &ps.Message{
// 						Data:       data,
// 						Attributes: publishInput.Attributes,
// 					}

// 					topicM.
// 						EXPECT().
// 						Publish(ctx, pubSubMsg).
// 						Return(resultM)

// 					resultM.
// 						EXPECT().
// 						Get(ctx).
// 						Return("", errors.New("Error to publish message"))

// 					errs = pubsubClient.Publish(ctx, publishInput)
// 					Expect(errs).To(Equal([]error{errors.New("Error to publish message")}))
// 				})
// 			})
// 		})

// 		Context("Success case", func() {
// 			When("publish message into pubsub topic successfully", func() {
// 				It("returns no errors", func() {
// 					data, _ := publishInput.Data.MarshalJSON()

// 					pubSubMsg := &ps.Message{
// 						Data:       data,
// 						Attributes: publishInput.Attributes,
// 					}

// 					topicM.
// 						EXPECT().
// 						Publish(ctx, pubSubMsg).
// 						Return(resultM)

// 					resultM.
// 						EXPECT().
// 						Get(ctx).
// 						Return("fake-server-id", nil)

// 					errs = pubsubClient.Publish(ctx, publishInput)
// 					Expect(errs).To(BeNil())
// 				})
// 			})
// 		})
// 	})
// })
