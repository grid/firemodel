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

protocol FiremodelDocumentSubscriber {
    associatedtype DocumentType
    func subscribe(receiver publish: @escaping (FiremodelDocumentEvent<DocumentType>) -> Void) -> Unsubscriber
}

enum FiremodelDocumentEvent<T> {
    case snapshot(_: T, metadata: SnapshotMetadata)
    case error(Error)
}

protocol FiremodelCollectionSubscriber {
    associatedtype DocumentType
    func subscribe(withQuery applyQuery: ((Query) -> Query)?,
                   receiver publish: @escaping (FiremodelCollectionEvent<DocumentType>) -> Void) -> Unsubscriber
}

enum FiremodelCollectionEvent<T> {
    case snapshot(_: [T], diff: (additions: [FiremodelChange<T>], modifications: [FiremodelChange<T>], removals: [FiremodelChange<T>]), metadata: SnapshotMetadata)
    case error(Error)
}

struct FiremodelChange<T> {
    let document: T
    let oldIndex: UInt
    let newIndex: UInt
}

// MARK: - Client

class FiremodelClient {
    private let firestore: FirebaseFirestore.Firestore

    static func dev() -> FiremodelClient {
        let firestore = FirebaseFirestore.Firestore.firestore()
        return FiremodelClient(firestore: firestore)
    }

    init(firestore: FirebaseFirestore.Firestore) {
        self.firestore = firestore
    }

    // MARK: - Root Collections

    func users() -> UserCollectionRef {
        return UserCollectionRef(ref: firestore.collection("users"), client: self)
    }

    func user(id: String) -> UserRef {
        return users().user(id: id)
    }

    // MARK: - Decoding

    func decode<T>(_ type: T.Type, from snapshot: FirebaseFirestore.DocumentSnapshot) throws -> T where T: Decodable {
        let decoder = DocumentSnapshotDecoder(documentSnapshot: snapshot,
                                              codingPath: [],
                                              userInfo: [firestoreClientDecodingKey: self])

        return try type.init(from: decoder)
    }

    func rawDocumentReference(_ path: String) -> DocumentReference {
        return self.firestore.document(path)
    }

    // MARK: - RPCs

//    private let transport: FiremodelRPCTransport
//
//    func dispatch(rpc: FiremodelRPCRequest,
//                  waitForUpdates: Bool,
//                  success: @escaping (FiremodelRPCSuccess) -> Void,
//                  failure: @escaping (FiremodelRPCError) -> Void) {
//
//        let meta = FiremodelRPCRequestMeta(bundleID: Bundle.main.bundleIdentifier,
//                                           appVersion: "1.0.0",
//                                           appBuild: "1",
//                                           osVersion: "13.0",
//                                           device: "iPhone4,3",
//                                           traceID: "TODO-trace",
//                                           spanID: "TODO-span")
//
//        let credentials = FiremodelRPCCredentials.firebaseAuthToken("myjwt")
//
//        try! self.transport.performRequest(
//            rpc,
//            meta: meta,
//            credentials: credentials,
//            success: { response in
//                var anyError: FiremodelRPCError?
//
//                let dispatchGroup = DispatchGroup()
//
//                if waitForUpdates {
//                    response.updatedDocuments.forEach({ updatedDoc in
//                        dispatchGroup.enter()
//                        updatedDoc.reference.getDocument(source: .server) { (snap, error) in
//                            if let error = error {
//                                debugPrint("Get document returned error")
//                                anyError = FiremodelRPCError.fromError(error)
//                            }
//                            debugPrint("Get document returned okk")
//                            dispatchGroup.leave()
//                        }
//                    })
//                }
//
//                dispatchGroup.notify(queue: DispatchQueue.main) {
//                    if let error = anyError {
//                        failure(error)
//                    } else {
//                        success(FiremodelRPCSuccess())
//                    }
//                }
//        },
//            failure: failure)
//    }
}

//enum FiremodelRPCCredentials {
//    case firebaseAuthToken(String)
//}
//
//protocol FiremodelRPCTransport {
//    func performRequest(_ rpc: FiremodelRPCRequest,
//                        meta: FiremodelRPCRequestMeta,
//                        credentials: FiremodelRPCCredentials,
//                        success: @escaping (FiremodelRPCResponse) -> Void,
//                        failure: @escaping (FiremodelRPCError) -> Void) throws
//}

//class FiremodelRPCHTTPTransport: FiremodelRPCTransport {
//    private let firestore: Firestore
//
//    private let session: URLSession
//
//    private let baseURL: URL
//
//    private let operationQueue: OperationQueue
//
//    init(firestore: FirebaseFirestore.Firestore, rpcBaseURL: URL) {
//        self.firestore = firestore
//        self.baseURL = rpcBaseURL
//        self.operationQueue = OperationQueue()
//        operationQueue.name = "FiremodelRPCHTTPTranspoprtQueue"
//        operationQueue.qualityOfService = .userInitiated
//
//        let config = URLSessionConfiguration()
//        config.allowsCellularAccess = true
//        config.httpShouldUsePipelining = true
//        config.waitsForConnectivity = true
//
//        session = URLSession(configuration: config,
//                             delegate: nil,
//                             delegateQueue: operationQueue)
//    }
//
//    private let encoder = JSONEncoder()
//    private let decoder = JSONDecoder()
//
//    func performRequest(_ rpc: FiremodelRPCRequest,
//                        meta: FiremodelRPCRequestMeta,
//                        credentials: FiremodelRPCCredentials,
//                        success: @escaping (FiremodelRPCResponse) -> Void,
//                        failure: @escaping (FiremodelRPCError) -> Void) throws {
//
//        guard let requestURL = URL(string: "/rpc", relativeTo: baseURL) else {
//            failure(FiremodelRPCError.unexpected("Invalid request url"))
//            return
//        }
//
//        var request = URLRequest(url: requestURL)
//        request.httpMethod = "POST"
//        request.addValue("application/json", forHTTPHeaderField: "content-type")
//        request.addValue("application/json", forHTTPHeaderField: "accept")
//        request.httpBody = try! encoder.encode(rpc)
//
//        session.dataTask(with: request) { [weak self] (data, response, error) in
//            if let error = error {
//                DispatchQueue.main.async {
//                    failure(.fromError(error))
//                }
//                return
//            }
//            guard let response = response as? HTTPURLResponse else {
//                DispatchQueue.main.async {
//                    failure(.unexpected("Missing HTTP response"))
//                }
//                return
//            }
//            guard response.mimeType == "application/json" else {
//                DispatchQueue.main.async {
//                    failure(.unexpected("Missing HTTP response"))
//                }
//                return
//            }
//            guard let data = data else {
//                DispatchQueue.main.async {
//                    failure(.unexpected("Missing HTTP response body"))
//                }
//                return
//            }
//            guard (200...300).contains(response.statusCode) else {
//                DispatchQueue.main.async {
//                    failure(.fromErrorResponse(response))
//                }
//                return
//            }
//
//            guard let decoder = self?.decoder, let firestore = self?.firestore else { return }
//            decoder.userInfo[firestoreClientDecodingKey] = firestore
//            let rpcResponse = try! decoder.decode(FiremodelRPCResponse.self, from: data)
//            DispatchQueue.main.async {
//                success(rpcResponse)
//            }
//        }.resume()
//    }
//}

let firestoreClientDecodingKey = CodingUserInfoKey(rawValue: "firestore")!

// RPCs

enum FiremodelRPCRequest {
    case sendMessage(to: FriendRef, content: Message)
}

struct FiremodelRPCRequestMeta: Codable {
    let bundleID: String?
    let appVersion: String?
    let appBuild: String?
    let osVersion: String?
    let device: String?
    let traceID: String?
    let spanID: String?
}


//extension FiremodelRPCRequest: Encodable {
//    func encode(to encoder: Encoder) throws {
//        var container = encoder.container(keyedBy: CodingKeys.self)
//        switch self {
//        case let .sendMessage(to: friendRef, content: messageContent):
//            try container.encode("SEND_MESSAGE", forKey: .type)
//            var nested = container.nestedContainer(keyedBy: SendMessageKeys.self, forKey: .sendMessage)
//            try nested.encode(friendRef, forKey: .to)
//            try nested.encode(messageContent, forKey: .content)
//        }
//    }
//
//    enum CodingKeys: CodingKey {
//        case type
//        case sendMessage
//    }
//
//    enum SendMessageKeys: CodingKey {
//        case to
//        case content
//    }
//}

extension FriendRef: Encodable {
    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(ref.path)
    }
}

//struct FiremodelRPCResponse: Decodable {
//    let updatedDocuments: [AnyDocumentReference]
//}
//
//// Wrapper around FirebaseFirestore.DocumentReference to facilitate extending the Firebase SDK class with Decodable.
//struct AnyDocumentReference {
//    fileprivate let reference: FirebaseFirestore.DocumentReference
//}
//
//extension AnyDocumentReference: Decodable {
//    init(from decoder: Decoder) throws {
//        guard let client = decoder.userInfo[firestoreClientDecodingKey] as? FiremodelClient else {
//            throw NSError(domain: "firemodel", code: 1, userInfo: nil)
//        }
//
//        let container = try decoder.singleValueContainer()
//        let rawReference = try container.decode(String.self)
//        self.reference = client.rawDocumentReference(rawReference)
//    }
//}
//
//struct FiremodelRPCError: Error {
//    let code: FiremodelRPCErrorCode
//    let message: String
//    let reason: String
//    let fieldErrors: [String: String]? = nil
//
//    static func fromError(_ error: Error) -> FiremodelRPCError {
//        return FiremodelRPCError(
//            code: .unknown,
//            message: error.localizedDescription,
//            reason: String(describing: error)
//        )
//    }
//
//    // Returns an error template for incorrect beahvior that should never happen.
//    static func unexpected(_ reason: String) -> FiremodelRPCError {
//        assertionFailure(reason)
//        return FiremodelRPCError(code: .unknown, message: "Something went wrong", reason: reason)
//    }
//
//    static func fromErrorResponse(_ response: HTTPURLResponse) -> FiremodelRPCError {
//        // TODO: Interpret status codes.
//        return FiremodelRPCError(code: .unknown,
//                                 message: "Something went wrong",
//                                 reason: "RPC server responded with \(response.statusCode) \(HTTPURLResponse.localizedString(forStatusCode: response.statusCode))")
//    }
//}
//
//enum FiremodelRPCErrorCode {
//    case unknown
//    case internalError
//    case badRequest
//    case preconditionFailed
//    case unimplemented
//    case unauthenticated
//}
//
//struct FiremodelRPCSuccess {
//}
//
// Models

struct UserCollectionRef {
    let ref: CollectionReference
    let client: FiremodelClient

    func user(id: String) -> UserRef {
        return UserRef(ref: ref.document(id), client: client)
    }
}

extension UserCollectionRef: FiremodelCollectionSubscriber {

    func subscribe(withQuery applyQuery: ((Query) -> Query)? = nil,
                   receiver publish: @escaping (FiremodelCollectionEvent<User>) -> Void) -> Unsubscriber {

        let registration = (applyQuery?(ref) ?? ref)
            .addSnapshotListener { (snap: QuerySnapshot?, error: Error?) in
                if let error = error {
                    publish(.error(error))
                    return
                }
                guard let snap = snap else {
                    assertionFailure("Error was nil but Snapshot was also nil. This is unexpected behavior from addSnapshotListener!")
                    publish(.error(FiremodelError.internalError))
                    return
                }


                var documents = [User]()
                var diff = (additions: [FiremodelChange<User>](), modifications: [FiremodelChange<User>](), removals: [FiremodelChange<User>]())
                for change in snap.documentChanges {
                    let user: User
                    do {
                        user = try self.client.decode(User.self, from: change.document)
                    } catch {
                        publish(.error(error))
                        return
                    }

                    documents.append(user)

                    let firemodelChange = FiremodelChange(document: user, oldIndex: change.oldIndex, newIndex: change.newIndex)

                    switch change.type {
                    case DocumentChangeType.added:
                        diff.additions.append(firemodelChange)
                    case DocumentChangeType.modified:
                        diff.modifications.append(firemodelChange)
                    case DocumentChangeType.removed:
                        diff.removals.append(firemodelChange)
                    default:
                        assertionFailure("unexpected firestore DocumentChangeType \(change.type)")
                    }

                }

                publish(.snapshot(documents, diff: diff, metadata: snap.metadata as SnapshotMetadata))
        }

        return Unsubscriber(listenerRegistration: registration)
    }
}

// MARK: - UserRef

struct UserRef {
    fileprivate let ref: DocumentReference
    fileprivate let client: FiremodelClient

    // MARK: Relations

    func grams() -> GramCollectionRef {
        return GramCollectionRef(ref: ref.collection("grams"), client: client)
    }

    func gram(id: String) -> GramRef {
        return GramRef(ref: ref.collection("grams").document(id), client: client)
    }

    func messages() -> MessageCollectionRef {
        return MessageCollectionRef(ref: ref.collection("messages"), client: client)
    }

    func message(id: String) -> MessageRef {
        return messages().message(id: id)
    }

    func friends() -> FriendCollectionRef {
        return FriendCollectionRef(ref: ref.collection("friends"), client: client)
    }

    func friend(id: String) -> FriendRef {
        return FriendRef(ref: ref.collection("friends").document(id), client: client)
    }
}

class Unsubscriber {
    private var listenerRegistration: ListenerRegistration?
    private var unsubscribeOnDeinit: Bool = true

    fileprivate init(listenerRegistration: ListenerRegistration) {
        self.listenerRegistration = listenerRegistration
    }

    // Shared prevents the automatic unsubscribe behavior for rare cases when the listener registration is never retained.
    func shared() {
        self.unsubscribeOnDeinit = false
    }

    deinit {
        if unsubscribeOnDeinit {
            unsubscribe()
        }
    }

    func unsubscribe() {
        listenerRegistration?.remove()
        listenerRegistration = nil
    }
}

enum FiremodelError: Error {
    case typeError
    case internalError
}

fileprivate extension User {
//    init(snapshot: DocumentSnapshot) throws {
//        guard let username: String? = snapshot.get("username") as? String? else {
//            throw FiremodelError.typeError
//        }
//        guard let displayName: String? = snapshot.get("displayName") as? String? else {
//            throw FiremodelError.typeError
//        }
//        guard let avatarUrl = try URL(snapshot: snapshot.get("avatar.url") as? String) else {
//            throw FiremodelError.typeError
//        }
//        guard let avatarColor = snapshot.get("avatar.color") as? String? else {
//            throw FiremodelError.typeError
//        }
//
//        self.init(username: username,
//                  displayName: displayName,
//                  avatar: Avatar(url: avatarUrl,
//                                 color: avatarColor))
//    }
}

extension User: Decodable {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        self.username = try container.decodeIfPresent(String.self, forKey: .username)
        self.displayName = try container.decodeIfPresent(String.self, forKey: .displayName)
        self.avatar = try container.decodeIfPresent(Avatar.self, forKey: .avatar)
    }

    enum CodingKeys: CodingKey {
        case username
        case displayName
        case avatar
    }
}

extension Avatar: Decodable {
    init(from decoder: Decoder) throws {
        let avatar = try decoder.container(keyedBy: CodingKeys.self)
        self.url = try avatar.decodeIfPresent(URL.self, forKey: .url)
        self.color = try avatar.decodeIfPresent(String.self, forKey: .color)
    }

    enum CodingKeys: CodingKey {
        case url
        case color
    }
}

extension Message: Encodable, Decodable {
    init(from decoder: Decoder) throws {
        let message = try decoder.container(keyedBy: CodingKeys.self)
        self.from = try message.decodeIfPresent(FriendRef.self, forKey: .from)
        let contentType = try message.decodeIfPresent(String.self, forKey: .content)
        let content = try message.nestedContainer(keyedBy: MessageContentType.self, forKey: .content)
        switch contentType {
        case MessageContentType.text.rawValue:
            let textMessageContent = try content.decode(TextMessageContent.self, forKey: .text)
            self.content = .text(textMessageContent)
        case MessageContentType.photo.rawValue:
            let photoMessageContent = try message.decode(PhotoMessageContent.self, forKey: .content)
            self.content = .photo(photoMessageContent)
        case let .some(contentType):
            self.content = .invalid(contentType)
        default:
            self.content = nil
        }
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        if let from = from {
            try container.encode(from, forKey: .from)
        } else {
            try container.encodeNil(forKey: .from)
        }
        var contentContainer = container.nestedContainer(keyedBy: MessageContentType.self, forKey: .content)
        if let content = content {
            switch content {
            case let .text(textMessageContent):
                try container.encode("TEXT", forKey: .content)
                try contentContainer.encode(textMessageContent, forKey: .text)
            case let .photo(photoMessageContent):
                try container.encode("PHOTO", forKey: .content)
                try contentContainer.encode(photoMessageContent, forKey: .photo)
            case let .invalid(invalid):
                try container.encode(invalid, forKey: .content)
            default:
                try container.encodeNil(forKey: .content)
            }
        }
    }

    private enum MessageContentType: String, CodingKey {
        case text
        case photo
    }

    private enum CodingKeys: CodingKey {
        case from
        case content
    }
}

extension TextMessageContent: Encodable, Decodable {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        self.message = try container.decode(String.self, forKey: .message)
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        if let message = message {
            try container.encode(message, forKey: .message)
        } else {
            try container.encodeNil(forKey: .message)
        }
    }

    enum CodingKeys: CodingKey {
        case message
    }
}

extension PhotoMessageContent: Encodable, Decodable {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        self.url = try container.decode(URL.self, forKey: .url)
        self.caption = try container.decode(String.self, forKey: .caption)
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(url, forKey: .url)
        try container.encode(caption, forKey: .caption)
    }

    enum CodingKeys: CodingKey {
        case caption
        case url
    }
}

extension FriendRef: Decodable {
    init(from decoder: Decoder) throws {
        guard let client = decoder.userInfo[firestoreClientDecodingKey] as? FiremodelClient else {
            assertionFailure("firestore client is missing in user info")
            throw DocumentSnapshotDecodingError.firestoreClientMissing
        }
        let container = try decoder.singleValueContainer()
        self.client = client
        self.ref  = client.rawDocumentReference(try container.decode(String.self))
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
//    init(snapshot: DocumentSnapshot) throws {
//        guard let rawSharedWith = snapshot.get("sharedWith") as? String? else {
//            throw FiremodelError.typeError(model: Gram.self, key: \Gram.sharedWith, expectedType: Audience.self, actualValue: snapshot.get("sharedWith"))
//        }
//        guard let photoUrl = try URL(snapshot: snapshot.get("photoUrl")) else {
//            throw FiremodelError.typeError(model: Gram.self, key: \Gram.photoUrl, expectedType: URL.self, actualValue: snapshot.get("photoUrl"))
//        }
//        guard let description = snapshot.get("description") as? String? else {
//            throw FiremodelError.typeError(model: Gram.self, key: \Gram.description, expectedType: String.self, actualValue: snapshot.get("description"))
//        }
//
//        self.init(sharedWith: Audience(snapshot: rawSharedWith),
//                  photoUrl: photoUrl,
//                  description: description)
//    }
}

fileprivate extension Message {
//    init(snapshot: DocumentSnapshot) throws {
//        guard let from = snapshot.get("from") as? DocumentReference? else {
//            throw FiremodelError.typeError
//        }
//
//        self.init(content: try MessageContent(snapshot: snapshot,
//                                              key: \Message.content,
//                                              fieldPath: ["messageContent"]),
//                  from: FriendRef(snapshot: from))
//    }
}

fileprivate extension FriendRef {
//    init?(snapshot ref: DocumentReference?) {
//        guard let ref = ref else { return nil }
//        self.init(ref: ref)
//    }
}

fileprivate extension MessageContent {
//    init?(snapshot: DocumentSnapshot, key: AnyKeyPath, fieldPath: [String]) throws {
//        let messageContent = snapshot.get(FieldPath(fieldPath)) as? String
//        switch messageContent {
//        case .none:
//            return nil
//        case .some("TEXT"):
//            let fp = FieldPath(fieldPath + ["message"])
//            let message = snapshot.get(fp) as? String
//            let textMessageContent = TextMessageContent(message: message)
//            self = .text(textMessageContent)
//        case .some("PHOTO"):
//            var fieldPath = fieldPath
//            fieldPath.append("\(fieldPath.removeLast()).photo")
//            self = .photo(try PhotoMessageContent(snapshot: snapshot, fieldPath: fieldPath))
//        case let .some(invalidMessageContent):
//            self = .invalid(invalidMessageContent)
//        }
//    }
}

//fileprivate extension Friend {
//    init(snapshot: DocumentSnapshot) throws {
//        let username = snapshot.get("username") as? String
//        let displayName = snapshot.get("displayName") as? String
//        let avatarUrl = try? URL(snapshot: snapshot.get(FieldPath(["avatar", "url"])
//        let avatarColor = snapshot.get(FieldPath(["avatar", "color"])) as? String
//        let friendsSince = snapshot.get("friendsSince") as? Timestamp
//
//        self.init(username: username,
//                  displayName: displayName,
//                  avatar: Avatar(url: avatarUrl,
//                                 color: avatarColor),
//                  friendsSince: friendsSince?.dateValue())
//    }
//}

//fileprivate extension TextMessageContent {
//    init(snapshot: DocumentSnapshot, fieldPath: [String]) throws {
//
//    }
//}

//fileprivate extension PhotoMessageContent {
//    init(snapshot: DocumentSnapshot, fieldPath: [String]) throws {
//        guard let caption = snapshot.get(FieldPath([ "caption"])) as? String? else {
//            throw FiremodelError.typeError
//        }
//        guard let url = try URL(snapshot: parent.get(FieldPath(["\(key).photo", "url"]))) else {
//            throw FiremodelError.typeError(model: <#T##Any.Type#>, key: <#T##AnyKeyPath#>, expectedType: <#T##Any.Type#>, actualValue: <#T##Any?#>)
//        }
//        self(caption: caption,
//             url: url)
//    }
//}

// MARK: - Rerferences

struct GramRef {
    fileprivate let ref: DocumentReference
    fileprivate let client: FiremodelClient

    func parent() -> UserCollectionRef {
        return UserCollectionRef(ref: ref.parent, client: client)
    }
}

struct GramCollectionRef {
    fileprivate let ref: CollectionReference
    fileprivate let client: FiremodelClient

    func parent() -> UserRef {
        return UserRef(ref: ref.parent!, client: client)
    }

    func gram(id: String) -> GramRef {
        return GramRef(ref: ref.document(id), client: client)
    }
}

struct MessageRef {
    fileprivate let ref: DocumentReference
    fileprivate let client: FiremodelClient

    func parent() -> MessageCollectionRef {
        return MessageCollectionRef(ref: ref.parent, client: client)
    }
}

struct MessageCollectionRef {
    fileprivate let ref: CollectionReference
    fileprivate let client: FiremodelClient

    func parent() -> UserRef {
        return UserRef(ref: ref.parent!, client: client)
    }

    func message(id: String) -> MessageRef {
        return MessageRef(ref: ref.document(id), client: client)
    }
}

struct FriendRef {
    fileprivate let ref: DocumentReference
    fileprivate let client: FiremodelClient

    func parent() -> UserCollectionRef {
        return UserCollectionRef(ref: ref.parent, client: client)
    }
}

struct FriendCollectionRef {
    fileprivate let ref: CollectionReference
    fileprivate let client: FiremodelClient

    func parent() -> UserRef {
        return UserRef(ref: ref.parent!, client: client)
    }

    func friend(id: String) -> FriendRef {
        return FriendRef(ref: ref.document(id), client: client)
    }
}
