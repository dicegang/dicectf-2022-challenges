const util = require("util");
const Multer = require("multer");
const MAXSIZE = 1024*1024; //1MB max upload

let processFile = Multer({
  storage: Multer.memoryStorage(),
  limits: { fileSize: MAXSIZE },
}).single("attack");

let processFileMiddleware = util.promisify(processFile);
module.exports = processFileMiddleware;

