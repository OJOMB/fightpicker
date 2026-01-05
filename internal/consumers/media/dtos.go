package media

/*
Example S3 Event Notification JSON (value of the Kafka record):
{
  "EventName": "s3:ObjectCreated:Put",
  "Key": "fightpicker-media-raw/users/019b38cb-8cf2-7ec5-897f-005b983dbd00/media/profile_image.jpeg",
  "Records": [
    {
      "eventVersion": "2.0",
      "eventSource": "minio:s3",
      "awsRegion": "",
      "eventTime": "2025-12-26T19:32:28.058Z",
      "eventName": "s3:ObjectCreated:Put",
      "userIdentity": {
        "principalId": "test_key"
      },
      "requestParameters": {
        "principalId": "test_key",
        "region": "",
        "sourceIPAddress": "10.89.4.9"
      },
      "responseElements": {
        "x-amz-id-2": "dd9025bab4ad464b049177c95eb6ebf374d3b3fd1af9251148b658df7ac2e3e8",
        "x-amz-request-id": "1884DAD1B4DD2793",
        "x-minio-deployment-id": "e1f19407-5835-4c0e-b168-e16dfe2deb1e",
        "x-minio-origin-endpoint": "http://10.89.4.9:9000"
      },
      "s3": {
        "s3SchemaVersion": "1.0",
        "configurationId": "Config",
        "bucket": {
          "name": "fightpicker-media-raw",
          "ownerIdentity": {
            "principalId": "test_key"
          },
          "arn": "arn:aws:s3:::fightpicker-media-raw"
        },
        "object": {
          "key": "users%2F019b38cb-8cf2-7ec5-897f-005b983dbd00%2Fmedia%2Fprofile_image.jpeg",
          "eTag": "d41d8cd98f00b204e9800998ecf8427e",
          "contentType": "image/jpeg",
          "userMetadata": {
            "content-type": "image/jpeg"
          },
          "sequencer": "1884DAD1B4F6B72E"
        }
      },
      "source": {
        "host": "10.89.4.9",
        "port": "",
        "userAgent": "PostmanRuntime/7.51.0"
      }
    }
  ]
}
*/

// EventAWS represents the structure of an S3 event notification
type EventAWS struct {
	EventVersion string           `json:"eventVersion"`
	EventSource  string           `json:"eventSource"`
	AWSRegion    string           `json:"awsRegion"`
	EventName    string           `json:"eventName"`
	EventTime    string           `json:"eventTime"`
	Records      []EventAWSRecord `json:"Records"`
}

type EventAWSRecord struct {
	S3 EventAWSRecordS3
}

type EventAWSRecordS3 struct {
	Bucket Bucket `json:"bucket"`
	Object Object `json:"object"`
}

// Bucket stores bucket related information from an S3 event
type Bucket struct {
	Name string `json:"name"`
	Arn  string `json:"arn"`
}

// Object stores object related information from an S3 event
type Object struct {
	Key         string `json:"key"`
	Size        int    `json:"size"`
	VersionId   string `json:"versionId"`
	ContentType string `json:"contentType"`
	ETag        string `json:"eTag"`
}
