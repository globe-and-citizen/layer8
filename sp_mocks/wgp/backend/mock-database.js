require("dotenv").config();
const BACKEND_URL = process.env.BACKEND_URL;

module.exports = {
    poems: [
            {
                "id": 1,
                "title": "The Red Wheelbarrow",
                "author": "WILLIAM CARLOS WILLIAMS",
                "body": "so much depends,\n upon \n a red wheel\nbarrow\nglazed with rain\nwater\nbeside the white\nchickens"
            },
            {   
                "id": 2,
                "title": "We Real Cool",
                "author": "Gwendolyn Brooks",
                "body": "We real cool. We\nLeft school. We\nLurk late. We\nStrike straight. We\nSing sin. We\nThin gin. We\nJazz June. We\nDie soon."
            },
            {
                "id": 3,
                "title": "The Road Not Taken",
                "author": "ROBERT FROST",
                "body": "Two roads diverged in a yellow wood,\nAnd sorry I could not travel both\nAnd be one traveler, long I stood\nAnd looked down one as far as I could\nTo where it bent in the undergrowth;"
            }
        ],
    users: [
        {
            "email": "star",
            "password": "$2b$10$ssDWc4sXNoafqEdAsvH8TOUXywGsFHsPEODTZlSB4AKe8cqe1PmCi",
            "profile_image": `${BACKEND_URL}/media/girl.png`,
        },
        {
            "email": "tester",
            "password": "$2b$10$vPCe/tNw/t2MHK/tGetY1exyvp4AhTC9w6mY5jyHHRAJrClfd1yYW",
            "profile_image": `${BACKEND_URL}/media/boy.png`,
        }
    ]
}