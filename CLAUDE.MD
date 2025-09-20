# Objective

Create a simplified version of a microblogging platform similar to Twitter that allows users to post, follow, and view a timeline of tweets.

# Requirements

## Tweets

* Users must be able to post short messages (tweets) that do not exceed a character limit (e.g., 280 characters).

## Follow

* Users must be able to follow other users.

## Timeline

* Users must be able to view a timeline displaying tweets from the users they follow.

# Assumptions

* All users are valid; there is no need to build a sign-in module or handle sessions. A user identifier can be sent via header, parameter, body, or whichever method is most convenient.
* The solution should be designed to scale to millions of users.
* The application must be optimized for read operations.

# Documentation

* Provide high-level documentation of the architecture and the components used in the application.
