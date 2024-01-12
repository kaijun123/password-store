## Password Storage

### Overview
This is a mini-project intended to show my understanding of the desired architecture of a system that stores user credentials, which involves the use of cryptographic salts and hashing.

In the database, the username, salt and hash is stored. The username needs to be unique. The salt is randomly generated. To obtain the hash, append the salt to the password and hash it using some hashing mechanism.

To check if the password provided is correct, retrieve the salt from the db for that username. Calculate the new hash and compare it with the stored hash. If they are equal, then the password is correct.

### References:
- bytebytego: https://www.linkedin.com/posts/alexxubyte_systemdesign-coding-interviewtips-activity-7114269933760307201-Bf80/
- Go and Gorm: https://dev.to/karanpratapsingh/connecting-to-postgresql-using-gorm-24fj
- Gin: https://gin-gonic.com/docs/quickstart/
- SHA256: https://gobyexample.com/sha256-hashes
- RNG: https://www.educative.io/answers/how-to-generate-random-numbers-in-go-language