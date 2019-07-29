//
//  Models.swift
//  Firemodel Demo
//
//  Created by Mickey Reiss on 7/27/19.
//  Copyright © 2019 Mickey Reiss. All rights reserved.
//

import Foundation

// MARK: - Structs

struct Avatar {
    let url: URL?
    let color: String?
}

// MARK: - Enums

enum Audience {
    case invalid(String)
    case global
    case friends
}

// MARKR: - Interfaces

protocol Userlike {
    var username: String? { get }
    var displayName: String? { get }
    var avatar: Avatar? { get }
}

// MARK: - Models

struct User: Userlike {
    let username: String?
    let displayName: String?
    let avatar: Avatar?
}

struct Gram {
    let sharedWith: Audience?
    let photoUrl: URL?
    let description: String?
}

struct Message {
    let content: MessageContent?
    let from: FriendRef?
}

enum MessageContent {
    case invalid(String)
    case text(String?)
    case photo(PhotoMessageContent)
}

struct PhotoMessageContent {
    let caption: String?
    let url: URL?
}

struct Friend: Userlike {
    let username: String?
    let displayName: String?
    let avatar: Avatar?
    let friendsSince: Date?
}

