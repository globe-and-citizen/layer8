const cors = require('cors');
const express = require('express');
const fs = require('fs');
const layer8 = require('layer8_middleware');

const app = express();
const port = 8000;
const upload = layer8.multipart({ dest: "uploads" });

app.use(layer8.tunnel);
app.use(cors());
app.use('/media/ex/', express.static('uploads'));
app.use('/media', layer8.static('uploads'));

app.get('/', (req, res) => {
  res.json({ message: "Hello there!" })
});

app.post("/api/upload", upload.single('file'), (req, res) => {
    const uploadedFile = req.file;

    if (!uploadedFile) {
        return res.status(400).json({ message: 'No file uploaded' });
    }
    
    res.status(200).json({
        message: "File uploaded successfully!",
        data: `${req.protocol}://${req.get('host')}/media/${req.file?.name}`
    });
});

app.get("/api/gallery", (req, res) => {
    const useExpress = req.query.use_express;
    var data = [];

    if (fs.existsSync("uploads")) {
        data = fs.readdirSync("uploads").map((file, i) => {
            return {
                id: i,
                name: file,
                url: `${req.protocol}://${req.get('host')}/media/${useExpress ? 'ex/' : ''}${file}`
            };
        });
    }

    res.status(200).json({ 
        message: "Your images are ready!",
        data: data
    });
});

app.listen(port, () => {
  console.log(`Image gallery app listening on port ${port}`)
})
