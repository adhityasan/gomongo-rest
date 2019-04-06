package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

// Gateway AWS API
type Gateway struct {
	Region    string
	KeyID     string
	SecretKey string
}

// CompareParam struct as a param for compare function
type CompareParam struct {
	ImgKTP    []byte
	ImgSelfie []byte
}

// Compare images using aws API
func (g *Gateway) Compare(p *CompareParam) (float64, error) {
	if sess, err := session.NewSession(&aws.Config{Region: aws.String(g.Region), Credentials: credentials.NewStaticCredentials(g.KeyID, g.SecretKey, "")}); err == nil {
		svcR := rekognition.New(sess)
		input := &rekognition.CompareFacesInput{
			SimilarityThreshold: aws.Float64(0),
			SourceImage: &rekognition.Image{
				Bytes: p.ImgKTP,
			},
			TargetImage: &rekognition.Image{
				Bytes: p.ImgSelfie,
			},
		}

		if res, err := svcR.CompareFaces(input); err == nil && len(res.FaceMatches) > 0 {
			for _, matchedFace := range res.FaceMatches {
				return *matchedFace.Similarity, nil
			}
		}
	}
	return 0, errors.New("Image comparison services is temporarily unable to process the request")
}

// Read text from images
func (g *Gateway) Read(img []byte) (string, error) {
	if sess, err := session.NewSession(&aws.Config{Region: aws.String(g.Region), Credentials: credentials.NewStaticCredentials(g.KeyID, g.SecretKey, "")}); err == nil {
		svcR := rekognition.New(sess)
		// Define input
		input := &rekognition.DetectTextInput{Image: &rekognition.Image{Bytes: img}}
		// Start Reading process
		if res, err := svcR.DetectText(input); err == nil && len(res.TextDetections) > 0 {
			var arr string
			for _, detectedtext := range res.TextDetections {
				if aws.StringValue(detectedtext.Type) == "LINE" {
					arr = arr + "\n" + aws.StringValue(detectedtext.DetectedText)
				}
			}
			return arr, nil
		}
	}
	return "", errors.New("OCR services is temporarily unable to process the request")
}
