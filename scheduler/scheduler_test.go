package scheduler

import (
	"sync"
	"testing"
	"time"

	"github.com/gerald-lbn/lazysinger/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/robfig/cron/v3"
)

func TestScheduler(t *testing.T) {
	log.SetLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scheduler Suite")
}

var _ = Describe("Scheduler", func() {
	var s *scheduler

	BeforeEach(func() {
		c := cron.New(cron.WithLogger(&logger{}))
		s = &scheduler{cron: c}
		s.cron.Start() // Start the scheduler for tests
	})

	AfterEach(func() {
		s.cron.Stop() // Stop the scheduler after tests
	})

	It("adds and executes a job", func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		executed := false
		id, err := s.AddJob("@every 100ms", func() {
			executed = true
			wg.Done()
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(id).ToNot(BeZero())

		wg.Wait()
		Expect(executed).To(BeTrue())
	})

	It("removes a job", func() {
		// Use a WaitGroup to ensure the job executes once
		wg := sync.WaitGroup{}
		wg.Add(1)

		counter := 0
		id, err := s.AddJob("@every 100ms", func() {
			counter++
			if counter == 1 {
				wg.Done() // Signal that the job has executed once
			}
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(id).ToNot(BeZero())

		// Wait for the job to execute at least once
		wg.Wait()

		// Verify job executed
		Expect(counter).To(Equal(1))

		// Remove the job
		s.RemoveJob(id)

		// Store the counter value
		currentCount := counter

		// Wait some time to ensure job doesn't execute again
		time.Sleep(200 * time.Millisecond)

		// Verify counter didn't increase
		Expect(counter).To(Equal(currentCount))
	})
})
