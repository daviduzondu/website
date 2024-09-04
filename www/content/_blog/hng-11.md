---
title: Surviving HNG 11
tags:
  - internship
  - hng
draft: false
date: 4 September 2024
---
If you are familiar with the Nigerian Tech Space, then you should have heard about HNG by now. It's a 2-month intensive program where thousands of people from around the world (mostly Nigerians) come together to put their skills to the test.

There are 10 stages in HNG. Each week, you are given a task that you must complete to move on to the next stage. In order to qualify as a finalist, you must make it to stage 10.

For HNG 11, I was added to a Slack Workspace alongside 16,000+ other humans. We were added to the `#stage-0` channel where we were given our first task. For stage 0, we the interns in the backend track were asked to write an article detailing a recent backend problem we solved. Our [submission](sk-experience.md) must contain at least two back links to HNG internship websites. The tasks in stage 0-1 were [interesting]() enough. I mean, the bar was low (for me), so I was able to breeze past them. 

In Stage 1, I had to set up a web server that exposes an API endpoint that will return the following JSON when a GET request is made to `/api/hello?vistor_name="Mark"`. 

```json
{
  "client_ip": "127.0.0.1", // The IP address of the requester
  "location": "New York", // The city of the requester
  "greeting": "Hello, Mark!, the temperature is 11 degrees Celcius in New York"
}
```

## Trouble Began In Stage 2 
Based on the [instructions](https://github.com/daviduzondu/hng-stage-2/blob/master/INSTRUCTIONS.md) for the Stage 2 task, we were supposed to setup PostgreSQL and connect our application to the database server. No, prior to this, I had never used PostgreSQL (or any other relational database).[^1] The only database system I was "experienced" in up to that point was MongoDB, because for some stupid reason I was scared of SQL. So I downloaded Postgres and installed it on my machine.

The instructions stated that we were allowed to use an ORM, but I decided to raw-dog SQL because I thought it'd be a great learning experience. For the database driver, I went with the [pg](https://npmjs.com/package/pg) package. From there things were a bit easier, I was familiar with Express so creating the API was not difficult at all. I hit another road block when it came to unit testing. Before Stage 2, I had never written tests for my code before. Ever! I just pushed my code to prod and pray to God that everything works. [^1]

I had to learn about unit testing overnight, and if you take a look a my code, it's pretty obvious that it was written by someone who had never written a unit test in his life. Either way, I score a 7/10 for this stage and I was allowed to proceed to stage 3.

## HNG is not as difficult as social media tells you
Trust me when I say this, the most difficult part of HNG is not coding. It's the collaboration/communication. I mean sure, if you are new to programming, you're probably NGMI, but if you know how to read and understand instructions, and you have good **communication** skills, you shouldn't have much of a problem going through the program. 

A lot of times, you'll be assigned a task with cryptic instructions, it is your job to decode what the instruction says as much as you can. The mentors are not always available (especially at the early stages of the program). Tens of thousands of people sign up for HNG every year, so you're competing with thousands of other confused interns for the attention of a mentor. 

Stage 3 was where chaos ensued. We were supposed to group ourselves (five interns per group) and work on an API and Database design. The instructions for this task was so confusing. Take a look at this: 

![](../../assets/stage-3.png)
API Design based on what? What the fuck does this even mean?

After a series of Google Meet calls and Huddles, We learnt that there was another repository containing a README file with better instructions. I have no clue why this was never mentioned in the original Slack message. Anyways, I worked with 4 other interns to deliver the API design/documentation and database blueprint. We passed, so we were moved to Stage 4. 

In Stage 4 we were supposed to create GitHub issues for the boilerplate repository we chose in the previous stage. Our issues would have to be approved before we can start working on them. In Stage 5, We were told that the Express repository was being discontinued and we should switch over to the Nest.js repository. Was I worried about deactivation because I didn't know shit about Nest? No. Did I lose my *steeze*? No.


![A screenshot from the HNG11 Slack Workspace](../../assets/Screenshot%20from%202024-09-04%2014-28-29.png)
The moment Mark Essien confirmed that the weeks of suffering and pain was now over.
[^1]: Shameful I know. 