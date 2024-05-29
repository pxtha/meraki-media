# Meradia

![flow.png](flow.png)

Link: https://merakilab-space.notion.site/Update-image-flow-c4d8b7ac4f4842049d7162621943db45?pvs=4

To upload an image to S3 storage using ReactJS and Golang, you can follow these steps:
### **Step 1: Call Backend to get Presigned URL**

- In the frontend (ReactJS), send an HTTP request to the backend (Golang) to get a presigned URL. Ensure that the backend is configured and has access to Amazon S3.
- The backend will receive the request, use its credentials to connect to Amazon S3, and create a presigned URL. This presigned URL is created with upload access (ACL) and has an expiration time.
- The backend will return the presigned URL to the frontend.

### **Step 2: Upload image to S3 using Presigned URL**

- In the frontend, use the received presigned URL to upload the image to S3. You can use libraries like Axios or Fetch to make an HTTP POST request to the presigned URL with the image data attached.
- When this request is sent, the image is transmitted from the frontend directly to S3 via the presigned URL.

### **Bước 3: Step 3: Confirm and save information to Backend**

- After the upload process is complete, the frontend can send another request to the backend to notify that the upload has been completed. This request may contain information about the uploaded image or any other related information.
- The backend can receive this request, confirm that the upload has been completed, and save the image information to the database or any other system.
- After successfully saving the information, the backend can return a link to the saved image, so the frontend can access and display the image.

This is a basic process to upload an image to S3 storage using ReactJS and Golang. However, note that the detailed setup of each step may vary depending on your requirements and configuration

---
# How to start source:
- Make sure you have read the [flow](flow.png) and understand it. By now you should understand how presigned url works and how to upload image to S3. Also you should know about imageproxy https://github.com/willnorris/imageproxy  and how to use it.
- Make sure you have the development environment matches with these notes below so we can mitigate any problems of version mismatch.
  - OS: Should use Linux (latest Ubuntu or your choice of distro) if possible. Windows does not play well with Docker and some other techs we may use. If you still prefer to use Windows, so you may have to cope with problems by yourself later since we're assuming everything will be developed and run on Linux.
  - Install Docker CE (latest) and docker-compose.
  - Install git for manage source code.
  - IDE of your choice, recommended Goland or VS Code.

### Development 
  - Clone code to local: ``` https://gitlab.com/merakilab9/meradia.git```
  - Add .env file: follow .env.example
  - Start development environment with Docker: ```make compose dev```
  - After started, services will be available at localhost with ports as below:
  ```Backend: 8099```
  ```ImageProxy: 8022```

### Testing api:

#### 1. PreUpload:
    ```curl
    curl --location 'localhost:8099/api/v1/media/pre-upload' \
    --header 'Content-Type: application/json' \
    --data '{
        "media_type": "png"
    }'
    ```

#### 2. Upload: BE test only; FE will use presigned url to upload
* get [PRESIGN_URL] from PreUpload api response
```curl
curl --location 'localhost:8099/api/v1/media/upload' \
--form 'upload_url="[PRESIGN_URL]"' \
--form 'file=@"/C:/Users/AlienWare/Downloads/3 copy 1.png"'
```

#### 3. Pos upload: Learn about imageproxy here: https://github.com/willnorris/imageproxy 
* get [FILE_NAME] from PreUpload api response
```curl
curl --location 'localhost:8099/api/v1/media/pos-upload' \
--header 'Content-Type: application/json' \
--data '{
    "name": "[FILE_NAME]",
    "media_type": "png",
    "type_to_crop": "avatar"
}'
```

#### Contact for support:
- Merakilab
- mail: pxthang97@gmail.com

---
# Understand the source code:
- The source code is written in Golang and base on Golang standard layout. You can read more about it here: https://github.com/golang-standards/project-layout
- The source code is divided into 3 main parts: handler, service, and repository. Each part has its own responsibility and should be separated to make the code clean and easy to maintain.
- Handler is responsible for handling HTTP requests, parsing request data, and returning responses. It should not contain business logic.
- Service is responsible for implementing business logic. It should not contain HTTP-related code.
- Repository is responsible for interacting with the database or external services. It should not contain business logic or HTTP-related code.
- These layer separation helps to make the code more modular, testable, and maintainable. It also makes it easier to change one part of the code without affecting other parts. They communicate with each other through interfaces, which allows for easy mocking and testing.