# Approval Service API Documentation

## Table of Contents

- [Approval Service API Documentation](#approval-service-api-documentation)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Authentication](#authentication)
  - [Base URL](#base-url)
  - [Endpoints](#endpoints)
    - [POST v1/approval/create](#post-v1approvalcreate)
    - [GET v1/approval/:id](#get-v1approvalid)
    - [GET /v1/approvals/:requestId](#get-v1approvalsrequestid)
    - [GET /v1/approvals?status=xxx?to=xxx?project=xxx?](#get-v1approvalsstatusxxxtoxxxprojectxxx)
    - [PUT v1/approval/update-status/:id](#put-v1approvalupdate-statusid)
    - [PUT v1/approval/sent-request/:id](#put-v1approvalsent-requestid)
    - [GET /v1/approval/user/:id?status=xxx\&to=xxx\&project=xx](#get-v1approvaluseridstatusxxxtoxxxprojectxx)
    - [GET /v1/approval/user/:id/receive-request?requsetUser=xx\&project=xx\&status=xx](#get-v1approvaluseridreceive-requestrequsetuserxxprojectxxstatusxx)
    - [GET /v1/approval/user/:id/send-request?to=xx\&project=xxx\&status=xx](#get-v1approvaluseridsend-requesttoxxprojectxxxstatusxx)
    - [DELETE /v1/approval/:id](#delete-v1approvalid)
  - [Status Codes](#status-codes)
  - [Event topics](#event-topics)

## Introduction

Welcome to the "Approval Service" API. Used to manage user Approval information in the system.

## Authentication

All requests to the API must include an API key in the header to authenticate:

`Authorization: Bearer YOUR_API_KEY`

## Base URL

All Approval service API endpoints are relative to this base
URL: `https://approval-service:{{app_port}}/`

## Endpoints

### POST v1/approval/create

worker create request approval.

**Headers**

```json
Authorization: JWT YOUR_TOKEN
```

**Request Body**:

```json
{
    "to":[1,2,3],
    "project": {"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"},
    "creation_date":1-10-5864,
    "request_user" :1,
    "task" :{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
}
```

**Response**:

```json
{
  "data" : {
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"pending",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":false
  },
  "message": "successful",
  "status": "ok",
  "status_code": 200
}
```

### GET v1/approval/:id

Get Approval by id

**Headers**

```json
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

```json
id : (require) type: uint
```

**Response**:

```json
{
  "data" : {
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"pending",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":false
  },
  "message": "successful",
  "status": "ok",
  "status_code": 200
}
```

### GET /v1/approvals/:requestId

Get the approvals by requestId
**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

- requestId (uuid)

**Response**:

```json
{

  "data" : [{
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"},
    "status" :"pending",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}],
    "IsSignature":false
  },
    {.....},// ... other profiles
  ],
  "status": "OK",
  "status_code": 200
}
```

### GET /v1/approvals?status=xxx?to=xxx?project=xxx?

Get All the approvals data. filter
**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

- to type: number
- project type: number
- status type: string

**Response**:

```json
{

  "data" : [{
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"pending",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":false
  },
    {.....},// ... other profiles
  ],
  "status": "OK",
  "status_code": 200
}
```

### PUT v1/approval/update-status/:id

Update status approval.

**Headers**

```json
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

```json
id : (require) type: uint
```

**Request Body:**

```json
{
  "status": "approved",
  "approver": 1,
  "is_signature": true,
}
```

**Response**:

```json
{
  "data" : {
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :1,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"approved",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":true
  },
  "message": "status changed",
  "status": "ok",
  "status_code": 200
}
```

### PUT v1/approval/sent-request/:id

teamlead request to HR or Approver.

**Headers**

```json
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

```json
id : (require) type: uint
```

**Request Body:**

```json
{
    "to":[1,2,3],
    "approver" :,
    "creation_date":10-10-5864,
    "request_user" :2,
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"HR",
    "is_signature":false
}
```

**Response**:

```json
{
  "data" : {
    "id" : 2,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f55e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"pending",
    "creationDate":1-10-5864,
    "requestUser" :2,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":false
  },
  "message": "status changed",
  "status": "ok",
  "status_code": 200
}
```

### GET /v1/approval/user/:id?status=xxx&to=xxx&project=xx

get approval from database by user id with filter

**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

- id : (require) type: uint
- requsetUser type: number
- project type: number
- status type: string

**Response**:

```json
{
  "data" : [{
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :1,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"approved",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":true
  },
    {.....},// ... other profiles
  ],
  "message": "successful",
  "status": "ok",
  "status_code": 200
}
```

### GET /v1/approval/user/:id/receive-request?requsetUser=xx&project=xx&status=xx

get approval from the database receives the user ID.

**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

- id : (require) type: uint
- requsetUser type: number
- project type: number
- status type: string

**Response**:

```json
{
  "data" : [{
   "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :1,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"}, //json
    "status" :"approved",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}], //json
    "IsSignature":true
  },
    {.....},// ... other profiles
  ],
  "message": "successful",
  "status": "ok",
  "status_code": 200
}
```

### GET /v1/approval/user/:id/send-request?to=xx&project=xxx&status=xx

get approval from the database sends the user ID.

**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

- id : (require) type: uint
- to type: number
- project type: number
- status type: string

**Response**:

```json
{
  "data" : [{
    "id" : 1,
    "requsetId":"185f6c4d-0b4e-4c1e-8d68-cbe862c9f38e",
    "name":"Timessheet for October 20 days",
	  "detail":"test test",
	  "name_request_user":"แทนไทย ทดสอบ",
	  "to_role":"teamlead",
    "to":[1,2,3],
    "approver" :1,
    "project" :{"id": 1,"name":"test","ข้อมูลอื่นๆ เพิ่มเติม"},
    "status" :"approved",
    "creationDate":1-10-5864,
    "requestUser" :1,
    "task" : [{"id": "1","ข้อมูลอื่นๆ เพิ่มเติม"},{....}],
    "IsSignature":true
  },
    {.....},// ... other profiles
  ],
  "message": "successful",
  "status": "ok",
  "status_code": 200
}
```

### DELETE /v1/approval/:id

delete approval from database

**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:

- id : (require) type: uint

**Response**:

```json
{
  "message": "approval deleted",
  "status": "ok",
  "status_code": 200
}
```

## Status Codes

<ul>
  <li>200 : OK. Request was successful.</li>
  <li>201 : Created. Resource was successfully created.</li>
  <li>400 : Bad request. The request was invalid or cannot be served.</li>
  <li>401 : Unauthorized. The request lacks valid authentication credentials.</li>
  <li>403 : Forbidden. The server understood the request, but it refuses to authorize it. This status code is similar to 401, but indicates that the client must authenticate itself to get permission.</li>
  <li>500 : Internal Server Error. The server has encountered a situation it does not know how to handle.</li>
</ul>

## Event topics

**ApprovalCreated**

**Publish** to the `tcchub-approval-approvalCreated` topic to publish information for create approval

**Message**:

```json
{
 "approval" : {
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :,
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [],
  },
}
```

**ApprovalUpdated**

**Publish** information about approval after update to the `tcchub-approval-approvalUpdated` topic

**Message**:

```json

{
 "approval" : {
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :,
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [],
  },
}
```

**ApprovalDeleted**

**Publish** delete event approval after delete to the `tcchub-approval-approvalDeleted` topic

**Message**:

```json
{
  "id":
}
```

API Documentation version:1 17/01/2024
