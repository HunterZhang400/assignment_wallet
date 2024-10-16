# assignment_wallet



# Install

- 1. Go to the project root, and then execute "docker compose build", and then "docker compose up", the demo services will be running.

- 2. Import ".postman_collection.json" to your postman, and you will see a collection name is "assignment_wallet".

- 3. Try APIs functions and enjoy it.




# Explains and Considerations

- 1. This business is about client's money, therefore, I need to add login and session control procedures.

- 2. I think not secure to operate client's account by a userID that send from API parameter, therefore the APIs only accept userID from session, that ensure all transactions operated by client in person.

- 3. To avoiding concurrent transaction conflicts involving same client, I adopted the distributed locker based on redis, and the locker granularity is every single client level, that could ensure data consistency but sacrificed little performance.

- 4. Other considerations: You may see secret configs in the etc/config.yaml, in real-world practice we should use complicate password and make it encrypted. The userID I suggest to used UUID string to make int unpredictable. Using bigint instead of money in postgresql due to its support not well. Added balance to response of deposit, withdraw and transfer business for better user experience, and clearly distinct the failure of deposit and query balance, avoiding make user confused to try again.

- 5. I encountered the issue of how to  mock redis and mysql gracefully and effectively, so the coverage not match your requirements, but I tested the all APIs by postman.

- 6. I use my free time to completing this assignment, about 4 hours design and conceived, 18 hours for coding, docker and gonglangci-lint etc.(PS, not the natural day hours, only working hours counted)



# How to review

## The complete call chain order is :
- Postman send a request
- router/router.go
- wallet/controller.go
- wallet/service.go
- pkg/db_util or pkg/redis_util
- 
 The most critical processes were located in wallet/service.go, which will be the primary focus area for reviewers.



This is my another assignment of implement to count inline comments and block comments in projects [assignment](https://github.com/HunterZhang400/go-homework-count-code-comments) , FYI.


