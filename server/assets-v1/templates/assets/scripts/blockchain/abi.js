export const abi = [
  {
    "inputs": [
      {
        "internalType": "address payable",
        "name": "_transactionAddress",
        "type": "address"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "constructor"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "BillAdded",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "BillPaid",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      }
    ],
    "name": "ClientCreated",
    "type": "event"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "internalType": "uint64",
        "name": "amount",
        "type": "uint64"
      },
      {
        "internalType": "uint64",
        "name": "timestamp",
        "type": "uint64"
      }
    ],
    "name": "addBillToClient",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "components": [
          {
            "internalType": "string",
            "name": "clientId",
            "type": "string"
          },
          {
            "internalType": "uint64",
            "name": "amount",
            "type": "uint64"
          },
          {
            "internalType": "uint64",
            "name": "timestamp",
            "type": "uint64"
          }
        ],
        "internalType": "struct PayAsYouGo.BillingInput[]",
        "name": "billings",
        "type": "tuple[]"
      }
    ],
    "name": "bulkAddBillToClient",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint64",
        "name": "rate",
        "type": "uint64"
      }
    ],
    "name": "changeAllClientRates",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "internalType": "uint64",
        "name": "rate",
        "type": "uint64"
      }
    ],
    "name": "changeRate",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address payable",
        "name": "_transactionAddress",
        "type": "address"
      }
    ],
    "name": "changeTransactionAddress",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "clientIDs",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "name": "clients",
    "outputs": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      },
      {
        "internalType": "uint256",
        "name": "unpaidBill",
        "type": "uint256"
      },
      {
        "internalType": "uint64",
        "name": "lastUsageFetchTime",
        "type": "uint64"
      },
      {
        "internalType": "uint64",
        "name": "rate",
        "type": "uint64"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      }
    ],
    "name": "getClientById",
    "outputs": [
      {
        "components": [
          {
            "internalType": "string",
            "name": "clientId",
            "type": "string"
          },
          {
            "internalType": "uint256",
            "name": "unpaidBill",
            "type": "uint256"
          },
          {
            "internalType": "uint64",
            "name": "lastUsageFetchTime",
            "type": "uint64"
          },
          {
            "internalType": "uint64",
            "name": "rate",
            "type": "uint64"
          },
          {
            "components": [
              {
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
              },
              {
                "internalType": "uint64",
                "name": "timestamp",
                "type": "uint64"
              },
              {
                "internalType": "enum PayAsYouGo.TransactionType",
                "name": "transactionType",
                "type": "uint8"
              }
            ],
            "internalType": "struct PayAsYouGo.Transaction[]",
            "name": "transactions",
            "type": "tuple[]"
          }
        ],
        "internalType": "struct PayAsYouGo.Client",
        "name": "",
        "type": "tuple"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getClients",
    "outputs": [
      {
        "components": [
          {
            "internalType": "string",
            "name": "clientId",
            "type": "string"
          },
          {
            "internalType": "uint256",
            "name": "unpaidBill",
            "type": "uint256"
          },
          {
            "internalType": "uint64",
            "name": "lastUsageFetchTime",
            "type": "uint64"
          },
          {
            "internalType": "uint64",
            "name": "rate",
            "type": "uint64"
          },
          {
            "components": [
              {
                "internalType": "uint256",
                "name": "amount",
                "type": "uint256"
              },
              {
                "internalType": "uint64",
                "name": "timestamp",
                "type": "uint64"
              },
              {
                "internalType": "enum PayAsYouGo.TransactionType",
                "name": "transactionType",
                "type": "uint8"
              }
            ],
            "internalType": "struct PayAsYouGo.Transaction[]",
            "name": "transactions",
            "type": "tuple[]"
          }
        ],
        "internalType": "struct PayAsYouGo.Client[]",
        "name": "",
        "type": "tuple[]"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint64",
        "name": "rate",
        "type": "uint64"
      },
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      }
    ],
    "name": "newClient",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "clientId",
        "type": "string"
      }
    ],
    "name": "payBill",
    "outputs": [],
    "stateMutability": "payable",
    "type": "function"
  }
]
