package main

import (
	"bytes"
	"fmt"

	"github.com/intermediate-service-ta/do"
	"github.com/intermediate-service-ta/gcs"
	"github.com/intermediate-service-ta/s3"
)

func main() {
	// GCS
	buf := new(bytes.Buffer)
	err := gcs.ListGCSBuckets(buf, "GOOG1EMGUURES77X6DIZEPF45TBM7ZH2ATQKX5PLUWB7FELW4DAN53G25T2II", "E/X5km0DcDDkab2k9XP1Czno98rPrF/T5LoOrXyM")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(buf.String())

	// Digital Ocean
	buf = new(bytes.Buffer)
	err = do.ListDOBuckets(buf)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(buf.String())

	// S3
	buf = new(bytes.Buffer)
	err = s3.ListS3Buckets(buf)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(buf.String())

}
