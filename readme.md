1. Create upload folder (mkdir upload)
2. Assumption all required GO packages already installed

To test upload:
1. Open index.html
2. Browse for valid image and click Submit.
3. upload.go will return JSON that original image title, shows url and resized_url
{"resized_url":"http://localhost/upload/doZPJTWUly1SIf3W_resized.jpg","title":"2 cor4.jpg","url":"http://localhost/upload/doZPJTWUly1SIf3W.jpg"}

To test download:
1. key in a valid image url
2. execute the function