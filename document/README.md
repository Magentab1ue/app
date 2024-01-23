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

Update create approval.

**Headers**

```json
Authorization: JWT YOUR_TOKEN
```

**Request Body**:

```json
{
    "to":,
    "project" :,
    "creationDate":,
    "requestUser" :,
    "task" :,
}
```

**Response**:
```json
{
  "data" : {
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :, 
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
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
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
  },
  "message": "successful",
  "status": "ok",
  "status_code": 200
}
```

### GET /v1/approvals?status=xxx?to=xxx?project=xxx?

Get the approvals data. filter
**Headers**:

```
Authorization: JWT YOUR_TOKEN
```

**Paremeters**:
- to  type: number
- project type: number
- status type: string

**Response**:

```json
{

  "data" : [{
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
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
  "status":,
  "approver":,
  "Approver":,
}
```

**Response**:

```json
{
  "data" : {
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
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
    "to":,
    "approver" :,
    "creationDate":,
    "requestUser" :,
    "is_signature":
}
```

**Response**:

```json
{
  "data" : {
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
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
- requsetUser  type: number
- project type: number
- status type: string

**Response**:

```json
{
  "data" : [{
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
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
- requsetUser  type: number
- project type: number
- status type: string

**Response**:

```json
{
  "data" : [{
    "id" : ,
    "requsetId":,
    "to":,
    "approver" :,
    "project" :, //json
    "status" :,
    "creationDate":,
    "requestUser" :,
    "task" : [], //json
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
- to  type: number
- project type: number
- status type: string

**Response**:

```json
{
  "data" : [{
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
