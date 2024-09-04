---
title: Surviving HNG 11
tags:
  - internship
  - hng
  - backend
draft: false
date: 4 September 2024
---
![A screenshot from the HNG11 Slack Workspace](../../assets/Screenshot%20from%202024-09-04%2014-28-29.png)
The moment Mark Essien confirmed that the weeks of suffering and pain was now over.

If you are familiar with the Nigerian Tech Space, then you should have heard about HNG by now. It's a 2-month intensive program where thousands of people from around the world (mostly Nigerians) come together to put their skills to the test.

There are 10 stages in HNG. Each week, you are given a task that you must complete to move on to the next stage. In order to qualify as a finalist, you must make it to stage 10.

For HNG 11, I was added to a Slack Workspace alongside 16,000+ other humans. We were added to the #stage-0 channel where we were given our first task. For stage 0, we the interns in the backend track were asked to write an article detailing a recent backend problem we solved. Our [submission](sk-experience.md) must contain at least two back links to HNG internship websites. The tasks in stage 0-1 were [interesting]() enough. I mean, the bar was low (for me), so I was able to breeze past them. 

In Stage 1, I had to set up a web server that exposes an API endpoint that will return the following JSON when a GET request is made to `/api/hello?vistor_name="Mark"`.

```json
{
  "client_ip": "127.0.0.1", // The IP address of the requester
  "location": "New York", // The city of the requester
  "greeting": "Hello, Mark!, the temperature is 11 degrees Celcius in New York"
}
```

## Trouble Began In Stage 2 
Based on the [instructions](https://github.com/daviduzondu/hng-stage-2/blob/master/INSTRUCTIONS.md) for the Stage 2 task, we were supposed to setup PostgreSQL and connect our application to the database server. No, prior to this, I had never used PostgreSQL (or any other relational database).[^1] The only database system I was "experienced" in up to that point was MongoDB, because for some stupid reason I was scared of SQL.

[^1]: Shameful I know. 