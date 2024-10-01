# spotify

# Rough overview of the Chatbot recommended project structure
project/
├── backend/
│   ├── main.go
│   ├── api/
│   │   ├── handlers.go
│   │   └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── model.go
│   └── ... (other backend-related files)
└── frontend/
    ├── public/
    │   └── index.html
    ├── src/
    │   ├── components/
    │   │   └── App.js
    │   ├── pages/
    │   │   └── Home.js
    │   ├── services/
    │   │   └── api.js
    │   ├── index.js
    │   └── ...
    ├── package.json
    ├── package-lock.json
    └── ... (other frontend-related files)