# Teach-Bot
Slack App for Education 

[Install on Slack](teach-bot.org)

##Runbook:
`make build` \
`make dev` 

At this point, you will be dropped into the container and a postgres container will be running as well -- run the following:

`go build` \
`./server`

##Current Features:

* **Acknowledge** - Slash command for producing a message that can track how many people have read a message and acknowledged it. Helpful for announcements.
* **Assignments**- Series of Slash commands and event triggers for helping a teacher assign a document (file or Google Drive) to students
* **Instructors** - Sets up a private conversation between a student and the instructors.
* **Anonymous Question** - Slash command that allows students to ask an anonymous question in a channel
* **Grades** - Slash Command that pretty prints a student's grades, for an instructor or for the student.