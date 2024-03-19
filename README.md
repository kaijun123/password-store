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
- Redis Key Expire: https://stackoverflow.com/questions/36172745/how-does-redis-expire-keys
- Session Management:
  - https://www.sohamkamani.com/golang/session-cookie-authentication/
  - https://www.sohamkamani.com/golang/password-authentication-and-storage/
  - https://github.com/alexedwards/scs
- Redis:
  - Redis Go Client: https://github.com/redis/go-redis?tab=readme-ov-file
  - Tutorial: https://tutorialedge.net/golang/go-redis-tutorial/
- Postman:
  - Cookies: https://learning.postman.com/docs/sending-requests/cookies/
- Gin:
  - https://techwasti.com/routing-in-go-gin-framework
  - https://medium.com/@ansujain/mastering-middleware-in-go-tips-tricks-and-real-world-use-cases-79215e72b4a8
  - https://stackoverflow.com/questions/62739044/how-to-share-variables-between-middlewares-in-go-gin
  - https://stackoverflow.com/questions/62608429/how-to-combine-group-of-routes-in-gin



https://blog.jetbrains.com/go/2022/11/22/comprehensive-guide-to-testing-in-go/
https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package
https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
https://medium.com/@gopal96685/handling-cookies-with-gin-framework-in-go-f119358c9cf3#:~:text=Setting%20a%20Cookie%20in%20Gin&text=You%20can%20use%20the%20SetCookie,The%20value%20of%20the%20cookie.
https://learning.postman.com/docs/sending-requests/cookies/#using-the-cookie-manager