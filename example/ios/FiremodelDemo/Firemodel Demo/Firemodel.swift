//
//  Firemodel.swift
//  Firemodel Demo
//
//  Created by Mickey Reiss on 7/26/19.
//  Copyright © 2019 Mickey Reiss. All rights reserved.
//

import Foundation
import FirebaseFirestore

typealias Source = FirebaseFirestore.FirestoreSource

// MARK: - Protocols

struct Change<T> {
    let document: T
    let oldIndex: Int
    let newIndex: Int
}

struct DocumentChange<T> {
    let document: T
    let oldIndex: UInt
    let newIndex: UInt
}

enum CollectionEvent<T> {
    case snapshot(_: [T], diff: (additions: [Change<T>], modifications: [Change<T>], removals: [Change<T>]), metadata: SnapshotMetadata)
    case error(Error)
}

protocol Subscribable {
    associatedtype T
    func subscribe(withQuery applyQuery: ((Query) -> Query)?,
                   receiver publish: @escaping (CollectionEvent<T>) -> Void) -> Unsubscriber
}

protocol Unsubscriber {
    func unsubscribe()
}

// MARK: - Client

class Firemodel {
    private let firestore: FirebaseFirestore.Firestore

    init() {
        firestore = FirebaseFirestore.Firestore.firestore()
    }

    // MARK: - Root Collections

    func users() -> UserCollectionRef {
        return UserCollectionRef(ref: firestore.collection("users"))

    }

    func user(id: String) -> UserRef {
        return users().user(id: id)
    }
}

struct UserCollectionRef {
    private let ref: CollectionReference

    init(ref: CollectionReference) {
        self.ref = ref
    }

    func user(id: String) -> UserRef {
        return UserRef(ref: ref.document(id))
    }
}

extension UserCollectionRef: Subscribable {
    typealias T = User

    func subscribe(withQuery applyQuery: ((Query) -> Query)? = nil,
                   receiver publish: @escaping (CollectionEvent<User>) -> Void) -> Unsubscriber {

        let registration = (applyQuery?(ref) ?? ref)
            .addSnapshotListener { (snap: QuerySnapshot?, error: Error?) in
                if let error = error {
                    publish(.error(error))
                    return
                }
                guard let snap = snap else {
                    assertionFailure("Error was nil but Snapshot was also nil. This is unexpected behavior from addSnapshotListener!")
                    return
                }


                var changes = [DocumentChange<User>]()
                for change in snap.documentChanges {

                    let user: User
                    do {
                        user = try User(snapshot: change.document)
                    } catch {
                        publish(.error(error))
                        return
                    }

                    changes.append(DocumentChange(document: user,
                                                   oldIndex: change.oldIndex,
                                                   newIndex: change.newIndex))
                }
                publish(.changes(changes, snap.metadata as SnapshotMetadata))
        }

        return ListenerRegistrationUnsubscriber(listenerRegistration: registration)
    }
}

// MARK: - UserRef

struct UserRef {
    private let ref: DocumentReference

    fileprivate init(ref: DocumentReference) {
        self.ref = ref
    }

    // MARK: Relations

    func grams() -> GramCollectionRef {
        return GramCollectionRef(ref: self.ref.collection("grams"))
    }

    func gram(id: String) -> GramRef {
        return GramRef(ref: ref.collection("grams").document(id))
    }

    func messages() -> MessageCollectionRef {
        return MessageCollectionRef(ref: ref.collection("messages"))
    }

    func message(id: String) -> MessageRef {
        return messages().message(id: id)
    }

    func friends() -> FriendCollectionRef {
        return FriendCollectionRef(ref: ref.collection("friends"))
    }

    func friend(id: String) -> FriendRef {
        return FriendRef(ref: ref.collection("friends").document(id))
    }

}


// MARK: UserRef Get Queries

extension UserRef: Gettable {

    func get(_ receiver: @escaping (Snapshot<User>) -> Void, source: Source = .default) {
        ref.getDocument(source: source) { (snap, error) in
            if let error = error {
                receiver(.error(error))
                return
            }
            guard let snap = snap else {
                // TODO
                return
            }

            do {
                let user = try User(snapshot: snap)
                receiver(.element(user))
            } catch {
                receiver(.error(error))
            }
        }
    }
}

// MARK: UserRef Watch Queries

struct UserSnapshot {
    let metadata: SnapshotMetadata
    let ref: UserRef
    let user: User
}

extension UserRef: Watchable {

    func watch(_ receiver: @escaping (Snapshot<User>) -> Void) -> Stopable {
        let remove = ref.addSnapshotListener { (snap, error) in
            if let error = error {
                receiver(.error(error))
                return
            }
            guard let snap = snap else {
                // TODO
                return
            }

            snap.metadata.timestamp

            do {
                let user = try User(snapshot: snap)
                receiver(.element(user))
            } catch {
                receiver(.error(error))
            }
        }
        return remove.stopper()
    }
}

private struct ListenerRegistrationUnsubscriber: Unsubscriber {
    private let listenerRegistration: ListenerRegistration

    func unsubscribe() {
        listenerRegistration.remove()
    }
}

enum FiremodelError: Error {
    case typeError
}

// User + Snapshot
fileprivate extension User {
    init(snapshot: DocumentSnapshot) throws {
        guard let username: String? = snapshot.get("username") as? String? else {
            throw FiremodelError.typeError
        }
        guard let displayName: String? = snapshot.get("displayName") as? String? else {
            throw FiremodelError.typeError
        }
        guard let avatarUrl = try URL(snapshot: snapshot.get("avatar.url") as? String) else {
            throw FiremodelError.typeError
        }
        guard let avatarColor = snapshot.get("avatar.color") as? String? else {
            throw FiremodelError.typeError
        }

        self.init(username: username,
                  displayName: displayName,
                  avatar: Avatar(url: avatarUrl,
                                 color: avatarColor))
    }
}

fileprivate extension Audience {
    init?(snapshot value: String?) {
        guard let value = value else { return nil }
        switch value {
        case "GLOBAL":
            self = .global
        case "FRIENDS":
            self = .friends
        default:
            self = .invalid(value)
        }
    }
}

fileprivate extension Gram {
    init(snapshot: DocumentSnapshot) throws {
        guard let rawSharedWith = snapshot.get("sharedWith") as? String? else {
            throw FiremodelError.typeError(model: Gram.self, key: \Gram.sharedWith, expectedType: Audience.self, actualValue: snapshot.get("sharedWith"))
        }
        guard let photoUrl = try URL(snapshot: snapshot.get("photoUrl")) else {
            throw FiremodelError.typeError(model: Gram.self, key: \Gram.photoUrl, expectedType: URL.self, actualValue: snapshot.get("photoUrl"))
        }
        guard let description = snapshot.get("description") as? String? else {
            throw FiremodelError.typeError(model: Gram.self, key: \Gram.description, expectedType: String.self, actualValue: snapshot.get("description"))
        }

        self.init(sharedWith: Audience(snapshot: rawSharedWith),
                  photoUrl: photoUrl,
                  description: description)
    }
}

fileprivate extension Message {
    init(snapshot: DocumentSnapshot) throws {
        guard let from = snapshot.get("from") as? DocumentReference? else {
            throw FiremodelError.typeError
        }

        self.init(content: try MessageContent(snapshot: snapshot,
                                              key: \Message.content,
                                              fieldPath: ["messageContent"]),
                  from: FriendRef(snapshot: from))
    }
}

fileprivate extension FriendRef {
    init?(snapshot ref: DocumentReference?) {
        guard let ref = ref else { return nil }
        self.init(ref: ref)
    }
}

fileprivate extension MessageContent {
    init?(snapshot: DocumentSnapshot, key: AnyKeyPath, fieldPath: [String]) throws {
        guard let messageContent = snapshot.get(FieldPath(fieldPath)) as? String? else {
            throw FiremodelError.typeError(model: MessageContent.self, key: key, expectedType: String.self, actualValue: snapshot.get("messageContent"))
        }

        switch messageContent {
        case .none:
            return nil
        case .some("TEXT"):
            var fieldPath = fieldPath
            fieldPath.append("\(fieldPath.removeLast()).text")
            guard let text = snapshot.get(fieldPath) as? String? else {
                throw FiremodelError.typeError(model: MessageContent.self, key: \MessageContent.self, expectedType: String.self, actualValue: snapshot.get(FieldPath(["\(key).text"])))
            }
            self = .text(text)
        case .some("PHOTO"):
            var fieldPath = fieldPath
            fieldPath.append("\(fieldPath.removeLast()).photo")
            self = .photo(try PhotoMessageContent(snapshot: snapshot, fieldPath: fieldPath))
        case let .some(messageContent):
            self = .invalid(messageContent)
        }
    }
}

fileprivate extension Friend {
    init(snapshot: DocumentSnapshot) throws {
        guard let username = snapshot.get("username") as? String? else {
            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
        }
        guard let displayName = snapshot.get("displayName") as? String? else {
            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
        }
        guard let avatarUrl = try URL(snapshot: snapshot.get(FieldPath(["avatar", "url"]))) else {
            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
        }
        guard let avatarColor = snapshot.get(FieldPath(["avatar", "color"])) as? String? else {
            throw FiremodelError.typeError(model: Avatar.self, key: \Avatar.color, expectedType: String.self, actualValue: snapshot.get("avatar.color"))
        }
        guard let friendsSince = snapshot.get("friendsSince") as? Timestamp? else {
            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
        }

        self.init(username: username,
                  displayName: displayName,
                  avatar: Avatar(url: avatarUrl,
                                 color: avatarColor),
                  friendsSince: friendsSince?.dateValue())
    }
}

fileprivate extension PhotoMessageContent {
    init(snapshot: DocumentSnapshot, fieldPath: [String]) throws {
        guard let caption = snapshot.get(FieldPath([ "caption"])) as? String? else {
            throw FiremodelError.typeError(model: PhotoMessageContent.self, key: \PhotoMessageContent.description, expectedType: String.self, actualValue: snapshot.get("description"))
        }
        guard let url = try URL(snapshot: parent.get(FieldPath(["\(key).photo", "url"]))) else {
            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
        }
        self(caption: caption,
             url: url)
    }
}

// MARK: - Rerferences

struct GramRef {
    private let ref: DocumentReference

    fileprivate init(ref: DocumentReference) {
        self.ref = ref
    }

    func parent() -> UserRef {
        return UserRef(ref: parentRef)
    }
}

struct GramCollectionRef {
    private let ref: CollectionReference
    private let parentRef: DocumentReference

    init(ref: CollectionReference, parent: DocumentReference) {
        self.ref = ref
    }

    func parent() -> UserRef {
        return UserRef(ref: parentRef)
    }

    func gram(id: String) -> GramRef {
        return GramRef(ref: ref.document(id))
    }
}

struct MessageRef {
    private let ref: DocumentReference

    fileprivate init(ref: DocumentReference) {
        self.ref = ref
    }

    func parent() -> MessageCollectionRef {
        return MessageCollectionRef(ref: ref.parent, parentRef: ref)
    }
}

struct MessageCollectionRef {
    private let ref: CollectionReference
    private let parentRef: DocumentReference

    fileprivate init(ref: CollectionReference, parentRef: DocumentReference) {
        self.ref = ref
        self.parentRef = parentRef
    }

    func parent() -> UserRef {
        return UserRef(ref: parentRef)
    }

    func message(id: String) -> MessageRef {
        return MessageRef(ref: ref.document(id))
    }
}

struct FriendRef {
    private let ref: DocumentReference

    fileprivate init(ref: DocumentReference) {
        self.ref = ref
    }

    func parent() -> UserCollectionRef {
        return UserCollectionRef(ref: ref.parent)
    }
}

struct FriendCollectionRef {
    private let ref: CollectionReference
    private let parentRef: DocumentReference

    fileprivate init(ref: CollectionReference, parentRef: DocumentReference) {
        self.ref = ref
    }

    func parent() -> UserRef {
        return UserRef(ref: parentRef)
    }

    func friend(id: String) -> FriendRef {
        return FriendRef(ref: ref.document(id))
    }
}

// MARK: - URL

fileprivate extension URL {
    init?(snapshot value: Any?) throws {
        switch value {
        case .none:
            return nil
        case let value as String:
            self.init(string: value)
        default:
            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
        }
    }
}

