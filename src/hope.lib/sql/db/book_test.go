package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"hope.lib/sql/db"
)

//http://onsi.github.io/ginkgo/

func DoTrue() bool {
	return true
}

var _ = Describe("Book", func() {
	var (
		longBook  db.User
		shortBook db.User
	)

	BeforeEach(func() {
		longBook = db.User{
			ID:  1,
			Name: "longBook",
			Level:  1488,
		}

		shortBook = db.User{
			ID:  2,
			Name: "shortBook",
			Level:  24,
		}
	})

	Describe("Check book desc", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.ID).To(Equal(int64(1)))
			})
		})

		Context("With fewer than 300 pages", func() {

			BeforeEach(func() {
				shortBook.Name = "test"
			})

			It("should be a short story", func() {
				Expect(shortBook.Name).To(Equal("test"))
			})
			It("panics in a goroutine", func(done Done) {
				go func() {
					defer GinkgoRecover()

					Ω(DoTrue()).Should(BeTrue())
					close(done)
				}()
			})
		})

		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(shortBook.Name).To(Equal("shortBook"))
			})
		})
	})
})
