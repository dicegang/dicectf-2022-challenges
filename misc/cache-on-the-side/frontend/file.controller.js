const processFile = require("./upload");
const { format } = require("util");
const { Storage } = require("@google-cloud/storage");
const { PubSub } = require("@google-cloud/pubsub");
const got = require("got");
const storage = new Storage();
const bucket = storage.bucket("cache-on-the-side");
const { v4: uuidv4 } = require('uuid');

const WAITFILE = "wait.html";
const STORAGE_BUCKET_URL = "https://storage.googleapis.com/cache-on-the-side/";

// Creates a client; cache this for further use
const topicName = "cache_challenge_entry";
const pubSubClient = new PubSub();

const createWaitingPage = async (req, res) => {
    // Pick new entry ID
    req.entry_id = uuidv4();

    const metadata = {
        contentType: 'text/html; charset=UTF-8',
        contentDisposition: 'inline',
        cacheControl: 'public, no-store',
    };

    await bucket.upload(WAITFILE, {
        destination: req.entry_id,
        metadata: metadata,
    });

    console.log(`Created waiting page ${req.entry_id}`);
}

const publishToQueue = async (req, res) => {
  const entry = {
    entry_id: req.entry_id,
    code_to_test: req.file.buffer.toString()
  }
  const data = Buffer.from(JSON.stringify(entry))

  const log_attributes = {
    ip: JSON.stringify(req.ip),
    ips: JSON.stringify(req.ips),
    x_real_ip: JSON.stringify(req.header("x-real-ip") ?? ""),
  }

  console.log(`Attributes: ${JSON.stringify(log_attributes)}`);
  console.log(`JSON buffer: ${data}`);

  const messageId = await pubSubClient.topic(topicName).publish(data, log_attributes);
  console.log(`Message ${messageId} published.`);
}


const upload = async (req, res) => {
  try {
    const recaptchaRes = await got({
      url: 'https://www.google.com/recaptcha/api/siteverify',
      method: 'POST',
      responseType: 'json',
      form: {
        secret: process.env.APP_RECAPTCHA_SECRET,
        response: req.body.recaptcha_code
      }
    })

    if (!recaptchaRes.body.success) {
      return res.status(200).send({message: "The reCAPTCHA is invalid"});
    }

    await processFile(req, res);

    if (!req.file) {
      return res.status(400).send({ message: "Please upload a file!" });
    }

    // Create the waiting page, a copy of "wait.html"
    await createWaitingPage(req, res);

    // Upload copy of file

    await publishToQueue(req, res);

    //redirect to wait.html
    res.redirect(301, STORAGE_BUCKET_URL+req.entry_id);

  } catch (err) {
    if (err.code == "LIMIT_FILE_SIZE") {
      return res.status(500).send({
        message: "File size cannot be larger than 1MB!",
      });
    }

    console.log(`Error! ${err}`);
    res.status(500).send({
      message: `Error. Please try again. ${err}`,
    });
  }
};

module.exports = {
  upload,
};
