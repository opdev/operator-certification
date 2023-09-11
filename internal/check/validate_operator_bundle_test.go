package check

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/opdev/operator-certification/internal/image"
)

var _ = Describe("BundleValidateCheck", func() {
	var bundleValidateCheck ValidateOperatorBundleCheck

	BeforeEach(func() {
		bundleValidateCheck = *NewValidateOperatorBundleCheck()
	})

	// AssertMetaData(&bundleValidateCheck)

	// TODO: Add more tests and bundles to testdata/ that exercise each of the
	// validations that we use.
	Describe("Operator Bundle Validate", func() {
		Context("When Operator Bundle Validate passes", func() {
			It("Should pass Validate", func() {
				imageRef := image.Reference{
					ImageFSPath: "./testdata/all_namespaces",
				}
				ok, err := bundleValidateCheck.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When Operator Bundle Validate does not Pass", func() {
			It("Should not pass Validate", func() {
				imageRef := image.Reference{
					ImageFSPath: "./testdata/invalid_bundle",
				}
				ok, err := bundleValidateCheck.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})
})
