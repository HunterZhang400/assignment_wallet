# assignment_wallet



# Install


# Explains and Considerations

## 1. This business is about client's money, therefore, I need to add login and session control procedures.

## 2. I think not secure to operate client's account by a userID that send from API parameter, therefore the APIs only accept userID from session, that ensure all transactions operated by client in person.

## 3. To avoiding concurrent transaction conflicts involving same client, I adopted the distributed locker based on redis, and the locker granularity is every single client level, that could ensure data consistency but sacrificed little performance.

## 4. Other considerations: You may see secret configs in the etc/config.yaml, in real-world practice we should use complicate password and make it encrypted. The userID I suggest to used UUID string to make int unpredictable. Using bigint instead of money in postgresql due to its support not well. Added balance to response of deposit, withdraw and transfer business for better user experience, and clearly distinct the failure of deposit and query balance, avoiding make user confused to try again.


# How to review

## The complete call chain order is :
- Postman send a request
- router/router.go
- wallet/controller.go
- wallet/service.go
## 

