//
//  Decoding.swift
//  Firemodel Demo
//
//  Created by Mickey Reiss on 8/6/19.
//  Copyright © 2019 Mickey Reiss. All rights reserved.
//

import Foundation
import FirebaseFirestore

struct DocumentSnapshotKey: CodingKey {
    let stringValue: String

    let intValue: Int? = nil

    init?(stringValue: String) {
        self.stringValue = stringValue
    }

    init?(intValue: Int) {
        return nil
    }
}

enum DocumentSnapshotDecodingError: Error {
    case firestoreClientMissing
}

struct DocumentSnapshotDecoder: Decoder {
    let documentSnapshot: DocumentSnapshot
    let codingPath: [CodingKey]
    let userInfo: [CodingUserInfoKey : Any]

    func container<Key>(keyedBy type: Key.Type) throws -> KeyedDecodingContainer<Key> where Key : CodingKey {
        return KeyedDecodingContainer(DocumentSnapshotKeyedDecodingContainerProtocol<Key>(documentSnapshot: documentSnapshot, codingPath: codingPath, userInfo: userInfo))
    }
    func singleValueContainer() throws -> SingleValueDecodingContainer {
        return DocumentSnapshotSingleValueDecodingContainer(documentSnapshot: documentSnapshot, codingPath: codingPath, userInfo: userInfo)
    }
    func unkeyedContainer() throws -> UnkeyedDecodingContainer {
        throw FiremodelError.typeError
    }
}

struct DocumentSnapshotUnkeyedDecodingContainer: UnkeyedDecodingContainer {
    let documentSnapshot: DocumentSnapshot

    let codingPath: [CodingKey]

    let userInfo: [CodingUserInfoKey: Any]

    var count: Int?

    var isAtEnd: Bool

    var currentIndex: Int

    mutating func decodeNil() throws -> Bool {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Bool.Type) throws -> Bool {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: String.Type) throws -> String {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Double.Type) throws -> Double {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Float.Type) throws -> Float {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int.Type) throws -> Int {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int8.Type) throws -> Int8 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int16.Type) throws -> Int16 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int32.Type) throws -> Int32 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: Int64.Type) throws -> Int64 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt.Type) throws -> UInt {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt8.Type) throws -> UInt8 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt16.Type) throws -> UInt16 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt32.Type) throws -> UInt32 {
        throw FiremodelError.internalError
    }

    mutating func decode(_ type: UInt64.Type) throws -> UInt64 {
        throw FiremodelError.internalError
    }

    mutating func decode<T>(_ type: T.Type) throws -> T where T : Decodable {
        throw FiremodelError.internalError
    }

    mutating func nestedContainer<NestedKey>(keyedBy type: NestedKey.Type) throws -> KeyedDecodingContainer<NestedKey> where NestedKey : CodingKey {
        throw FiremodelError.internalError
    }

    mutating func nestedUnkeyedContainer() throws -> UnkeyedDecodingContainer {
        throw FiremodelError.internalError
    }

    mutating func superDecoder() throws -> Decoder {
        throw FiremodelError.internalError
    }
}

struct DocumentSnapshotSingleValueDecodingContainer: SingleValueDecodingContainer {
    let documentSnapshot: DocumentSnapshot
    let codingPath: [CodingKey]
    let userInfo: [CodingUserInfoKey: Any]

    func decodeNil() -> Bool {
        return false
    }

    func decode(_ type: Bool.Type) throws -> Bool {
        throw FiremodelError.internalError
    }

    func decode(_ type: String.Type) throws -> String {
        throw FiremodelError.internalError
    }

    func decode(_ type: Double.Type) throws -> Double {
        throw FiremodelError.internalError
    }

    func decode(_ type: Float.Type) throws -> Float {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int.Type) throws -> Int {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int8.Type) throws -> Int8 {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int16.Type) throws -> Int16 {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int32.Type) throws -> Int32 {
        throw FiremodelError.internalError
    }

    func decode(_ type: Int64.Type) throws -> Int64 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt.Type) throws -> UInt {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt8.Type) throws -> UInt8 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt16.Type) throws -> UInt16 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt32.Type) throws -> UInt32 {
        throw FiremodelError.internalError
    }

    func decode(_ type: UInt64.Type) throws -> UInt64 {
        throw FiremodelError.internalError
    }

    func decode<T>(_ type: T.Type) throws -> T where T : Decodable {
        throw FiremodelError.internalError
    }
}

struct DocumentSnapshotKeyedDecodingContainerProtocol<Key>: KeyedDecodingContainerProtocol where Key : CodingKey {
    let documentSnapshot: DocumentSnapshot
    let codingPath: [CodingKey]
    let userInfo: [CodingUserInfoKey: Any]

    var allKeys: [Key] {
        return Array(self.documentSnapshot.data()?.keys.compactMap { Key(stringValue: $0) } ?? [])
    }

    func contains(_ key: Key) -> Bool {
        return self.documentSnapshot.data()?.keys.contains(key.stringValue) ?? false
    }

    func decodeNil(forKey key: Key) throws -> Bool {
        return documentSnapshot.get(key.stringValue) == nil
    }

    private func primativeValue(forKey key: Key) throws -> Any {
        let fp = FieldPath(codingPath.map { $0.stringValue } + [key.stringValue])
        guard case let .some(value) = documentSnapshot.get(fp) else {
            throw DecodingError.keyNotFound(key, DecodingError.Context(codingPath: codingPath, debugDescription: "key missing"))
        }
        return value
    }

    func decode(_ type: Bool.Type, forKey key: Key) throws -> Bool {
        guard let value = try primativeValue(forKey: key) as? Bool else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: String.Type, forKey key: Key) throws -> String {
        guard let value = try primativeValue(forKey: key) as? String else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Double.Type, forKey key: Key) throws -> Double {
        guard let value = try primativeValue(forKey: key) as? Double else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Float.Type, forKey key: Key) throws -> Float {
        guard let value = try primativeValue(forKey: key) as? Float else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int.Type, forKey key: Key) throws -> Int {
        guard let value = try primativeValue(forKey: key) as? Int else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int8.Type, forKey key: Key) throws -> Int8 {
        guard let value = try primativeValue(forKey: key) as? Int8 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int16.Type, forKey key: Key) throws -> Int16 {
        guard let value = try primativeValue(forKey: key) as? Int16 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int32.Type, forKey key: Key) throws -> Int32 {
        guard let value = try primativeValue(forKey: key) as? Int32 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: Int64.Type, forKey key: Key) throws -> Int64 {
        guard let value = try primativeValue(forKey: key) as? Int64 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt.Type, forKey key: Key) throws -> UInt {
        guard let value = try primativeValue(forKey: key) as? UInt else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt8.Type, forKey key: Key) throws -> UInt8 {
        guard let value = try primativeValue(forKey: key) as? UInt8 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt16.Type, forKey key: Key) throws -> UInt16 {
        guard let value = try primativeValue(forKey: key) as? UInt16 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt32.Type, forKey key: Key) throws -> UInt32 {
        guard let value = try primativeValue(forKey: key) as? UInt32 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode(_ type: UInt64.Type, forKey key: Key) throws -> UInt64 {
        guard let value = try primativeValue(forKey: key) as? UInt64 else {
            throw DecodingError.typeMismatch(type, DecodingError.Context(codingPath: codingPath, debugDescription: "invalid type"))
        }
        return value
    }

    func decode<T>(_ type: T.Type, forKey key: Key) throws -> T where T : Decodable {
        let t = DocumentSnapshotDecoder(documentSnapshot: documentSnapshot, codingPath: codingPath + [key], userInfo: self.userInfo)
        return try type.init(from: t)
    }

    func nestedContainer<NestedKey>(keyedBy type: NestedKey.Type, forKey key: Key) throws -> KeyedDecodingContainer<NestedKey> where NestedKey : CodingKey {
        return KeyedDecodingContainer(DocumentSnapshotKeyedDecodingContainerProtocol<NestedKey>(documentSnapshot: documentSnapshot, codingPath: codingPath + [key], userInfo: userInfo))
    }

    func nestedUnkeyedContainer(forKey key: Key) throws -> UnkeyedDecodingContainer {
        guard let value = try primativeValue(forKey: key) as? [Any] else {
            throw DecodingError.typeMismatch([Any].self, DecodingError.Context(codingPath: codingPath, debugDescription: "Unexpected type for key \(key)"))
        }
        return DocumentSnapshotUnkeyedDecodingContainer(documentSnapshot: documentSnapshot, codingPath: codingPath + [key], userInfo: userInfo, count: value.count, isAtEnd: value.isEmpty, currentIndex: 0)
    }

    func superDecoder() throws -> Decoder {
        throw FiremodelError.internalError
    }

    func superDecoder(forKey key: Key) throws -> Decoder {
        throw FiremodelError.internalError
    }
}


