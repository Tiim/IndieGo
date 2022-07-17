# Go Comment API

This project is a simple but extendable api to store comments for a jamstack website.

It supports prerendered apps by returning all comments on build time, and only returning comments since the last build when dynamically queried.

> [Blogpost why I built this project](https://tiim.ch/blog/2022-07-12-first-go-project-commenting-api)

![Image of the Go gopher with a speech bubble](/go-comment-api-image.svg)

## Demo

You can try out this project on my [Blogpost about this project](https://tiim.ch/blog/2022-07-12-first-go-project-commenting-api).

## API

The website can interact with the comment api via a rest(-ish) API. All routes are described here.

### Fetch all comments

This endpoint returns all comments stored in the database. This can be usefull when pre-rendering your app, so you only have to make a single api call.

```http
GET /comment
```
**Returns**

```json
[
  {
    "id":"c5e9ecff-1698-452b-880a-d6f0d934e740",
    "reply_to":"0791f666-b76b-455d-b543-0c4fe9fb7a44",
    "timestamp":"2022-07-11T12:17:32Z",
    "page":"blog/my-blogpost-slug",
    "content":"This is a test comment!",
    "name":"John"
  }
]
```

#### Query Parameters

```querystring
?since=(timestamp)
```

- **timestamp**: (Optional) is a ISO8602 formatted timestamp. The results will be filtered to only contain comments created after this timestamp.

### Fetch comments for one page

This endpoint returns all comments for a single page. This is the recommended way to load comments in a jam-stack website.

```http
GET /comment/:page
```
- **page**: The url encoded identifier of the page. This could be a page slug, the full url or an arbitrary id like a uuid.

**Returns** See [Fetch all comments](#fetch-all-comments).

#### Query Parameters

See [Fetch all comments](#fetch-all-comments).

### New comment

This is the API endpoint used to create a new comment from your comment form.

```http
POST /comment
```
**Body**
```json
{
  "name": "Jane",
  "page": "blog/my-blogpost-slug",
  "content": "This is a nice commenting api!",
  "reply_to": "0791f666-b76b-455d-b543-0c4fe9fb7a44",
  "email": "email@example.com",
  "notify": true
}
```
- **name** (Required) the name of the commenter.
- **page** (Required) the page id. Can be for example the slug of the page, the full url or a uuid.
- **content** (Required) the comment itself.
- **reply_to** (Optional) the id of another comment.
- **email** (Optional) the email of the commenter. Shown in the [Dashboard](#dashboard) and used to send reply notifications.
- **notify** (Optional) opt-in parameter for reply notifications. By default it is false. Always let the commenter decide if they want to opt in!

**Returns** a single comment in the format as seen in [Fetch all comments](#fetch-all-comments).

## User interface

The api currently has a few very basic user interfaces: The dashboard, a comment unsubscribe page and an email unsubscribe page.

### Dashboard

The dashboard shows a list of all comments and allows you to delete them.

The url for the dashboard is `/admin`. The dasboard requires authentication with a username and password. The username is "admin" and the password can be specified via the env variable `ADMIN_PW`.

### Unsubscribe from comment

This page is meant to be shown when clicking the unsubscribe link from a notification email. The url is `/unsubscribe/comment/:unsubscribe-secret`. Opening this page with a valid unsubscribe secret will mark the comment with `notify=false`.

### Unsubscribe from all emails

This page is meant to be shown when clicking the unsubscribe all link from the previous page. The url is `/unsubscribe/email/:email-address`. Opening this page with an email will mark all comments that use this exact email address with `notify=false`.

## Customization

The api is built as modular pieces which you can combine and extend as you wish. The components are:
- the API component
- the model component
- the event component

The api component handles all the api routes. It is not meant to be extended directly. 

The model component handles the storage of the comments. Currently the implementation of the model uses a SQLite database. You can add support for different data stores like other data base engines or a file based store by implementing the `Store` and `SubscriptionStore` interfaces defined in [store.go](/model/store.go).

The event component is an easy way to hook into what happens when a new comment gets submitted and when a comment gets deleted.
By default there are two event components: the [EmailNotify](/event/emailnotify.go) and the [ReplyEmailNotify](/event/replyemailnotify.go) components. You can add as many event handlers as you like.
