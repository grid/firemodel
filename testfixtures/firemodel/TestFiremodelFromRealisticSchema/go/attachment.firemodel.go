// DO NOT EDIT - Code generated by firemodel (dev).

package firemodel

// Firestore document location: /users/{user_id}/attachments/{attachment_id}
type Attachment struct {
	Title   string            `firestore:"title,omitempty"`
	Content AttachmentContent `firestore:"content,omitempty"`
}