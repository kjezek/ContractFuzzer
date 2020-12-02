/*
    This is a server that orchestrates experiments
 */

const express = require("express");
const fs = require("fs");

const app = express();

function readTasks(inputDir) {
    // read all available tasks
    const dir = fs.opendirSync(inputDir)
    let tasks = []
    let dirent;

    while ((dirent = dir.readSync()) !== null) {
        tasks.push(parseInt(dirent.name));
    }

    dir.closeSync();
    return tasks.sort((a, b)=> a-b);
}

function server() {

    const tasks = readTasks(inputDir);
    let index = 0;

    app.listen(9999, () => {
        console.log("Server running on port 9999");
    });

    // Return next available task
    app.get("/task", (req, res, next) => {
        const nextTask = index === tasks.length ? "DONE" : tasks[index++]
        console.log("Next task is " + nextTask + " Index: " + index + "/" + tasks.length)
        res.send(nextTask.toString());
    });

}

const args = process.argv.slice(2);
const inputDir = args[0];

server();

