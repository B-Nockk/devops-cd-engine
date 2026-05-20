bucket         = "YOUR_S3_BUCKET_NAME"
key            = "terraform/state/cd-engine.tfstate"
region         = "us-east-1"
dynamodb_table = "YOUR_DYNAMODB_TABLE_NAME"   # Optional, for state locking
encrypt        = true
