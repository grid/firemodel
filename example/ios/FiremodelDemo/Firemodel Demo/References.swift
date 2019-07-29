//
//  References.swift
//  Firemodel Demo
//
//  Created by Mickey Reiss on 7/27/19.
//  Copyright © 2019 Mickey Reiss. All rights reserved.
//

import Foundation
import FirebaseFirestore

struct GramRef {
    private let ref: DocumentReference

    fileprivate init(ref: DocumentReference) {
        self.ref = ref
    }
}

struct GramCollectionRef {
    private let ref: CollectionReference

    init(ref: CollectionReference) {
        self.ref = ref
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
}

struct MessageCollectionRef {
    private let ref: CollectionReference

    fileprivate init(ref: CollectionReference) {
        self.ref = ref
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
}

struct FriendCollectionRef {
    private let ref: CollectionReference

    fileprivate init(ref: CollectionReference) {
        self.ref = ref
    }

    func friend(id: String) -> FriendRef {
        return FriendRef(ref: ref.document(id))
    }
}
