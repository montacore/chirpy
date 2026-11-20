#!/usr/bin/env bash

#TEST: JSON string body too long 
# Should be expecting error code: 400
# and error body
curl -i -X POST http://localhost:8080/api/validate_chirp \
    -H "Content-Type: application/json" \
    -d '{"body": "Spent the whole morning repainting the fence, only to realize halfway through that the color I picked looks nothing like the sample card at the hardware store. Lesson learned."}'

#TEST: Valid(should get 200 response)
curl -i -X POST http://localhost:8080/api/validate_chirp \
    -H "Content-Type: application/json" \
    -d '{"body": "Fresh coffee and a quiet morining"}'

#TEST: Broken JSON(tests decode error branch)
curl -i -X POST http://localhost:8080/api/validate_chirp \
  -H "Content-Type: application/json" \
  -d '{"body": oops}'

#TEST: Valid JSON with "profanity" that should be returned as cleaned
curl -i -X POST http://localhost:8080/api/validate_chirp \
    -H "Content-Type: application/json" \
    -d '{"body": "Wow, quite the kerfuffle that was. what is a fornax sharbert"}'
