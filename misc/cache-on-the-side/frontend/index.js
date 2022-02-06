const express = require("express");
const app = express();
const router = express.Router();
const controller = require("./file.controller");

let initRoutes = (app) => {
    router.post("/upload", controller.upload);
    app.use(router);
  };

app.set('views', './public');
app.set('view engine', 'ejs')

app.use(express.urlencoded({ extended: false }));
initRoutes(app);

app.get("/", (req, res) => {
    console.log(' Return /');
    res.render('index');
});

app.use((err, req, res, next) => {
    console.log(err);
    res.redirect("/");
});

const port = process.env.PORT || 80;
app.listen(port, () => {
  console.log(`Running at localhost:${port}`);
});
