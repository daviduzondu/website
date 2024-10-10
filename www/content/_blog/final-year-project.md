---
title: Creating an Online Test Management System
draft: true
---


- Went with Firebase Auth - Could not figure out how to sync user data with backend. Thinking about switching to Clerk
- Thinking about creating a global guard in the Nest.js backend that checks if the user data from Firebase is in sync with user data in the DB.  
- I love Firebase, I don't want to switch to Clerk. I'm not familiar with Clerk.
- `getRedirectResult` keeps returning null
- Realized that since I am using Just Firebase Auth, I am actually investing more time into making my app work with just Firebase Auth. I already have the major FE components for authentication. I should just work with Passport.js at this point. 
- Learning about Passport.js strategies LMAO
- I spent hours pulling my hair wondering why nullable columns aren't shown with `$inferInsert`, turns out the solution to the problem is [here](https://github.com/drizzle-team/drizzle-orm/issues/2636#issuecomment-2314537568)
- Switched from Drizzle to Kysesly. It feels good to be home. Honestly.
- Implementing the Equation Editor was such a pain in the ass. I generated most of the UI with [https://v0.dev](v0.dev)
Resources
TestGorilla
https://blog.hyperknot.com/p/comparing-auth-providers
https://stackoverflow.com/questions/70170262/synchronize-users-created-with-firebase-auth-to-my-custom-backend/70171351#70171351:~:text=Problem%3A%20What%20if%20my%20backend%20call%20fails%20but%20the%20register%20process%20was%20successful

