// DO NOT EDIT - Code generated by firemodel (dev).

import Foundation
import Pring

// TODO: Add documentation to TestEnum in firemodel schema.
@objc enum TestEnum: Int {
    // TODO: Add documentation to Left in firemodel schema.
    case left
    // TODO: Add documentation to Right in firemodel schema.
    case right
    // TODO: Add documentation to Up in firemodel schema.
    case up
    // TODO: Add documentation to Down in firemodel schema.
    case down
}

extension TestEnum: CustomDebugStringConvertible {
    init?(firestoreValue value: Any?) {
        guard let value = value as? String else {
            return nil
        }
        switch value {
        case "LEFT":
            self = .left
        case "RIGHT":
            self = .right
        case "UP":
            self = .up
        case "DOWN":
            self = .down
        default:
            return nil
        }
    }

    var firestoreValue: String? {
        switch self {
        case .left:
            return "LEFT"
        case .right:
            return "RIGHT"
        case .up:
            return "UP"
        case .down:
            return "DOWN"
        }
    }

    var debugDescription: String { return firestoreValue ?? "<INVALID>" }
}

// TODO: Add documentation to TestStruct in firemodel schema.
@objcMembers class TestStruct: Pring.Object {
  // TODO: Add documentation to where in firemodel schema.
  var where: String?
  // TODO: Add documentation to how_much in firemodel schema.
  var howMuch: Int = 0
}

// A Test is a test model.
@objcMembers class TestModel: Pring.Object {
    static var userId: String = ""
static var testModelId: String = ""
    override class var path: String { return "users/\(userId)/test_models/\(testModelId)" }
    // The name.
    dynamic var name: String?
    // The age.
    dynamic var age: Int = 0
    // The number pi.
    dynamic var pi: Float = 0.0
    // The birth date.
    dynamic var birthdate: Date?
    // True if it is good.
    dynamic var isGood: Bool = false
    // TODO: Add documentation to data in firemodel schema.
    dynamic var data: Data?
    // TODO: Add documentation to friend in firemodel schema.
    dynamic var friend: Pring.Reference<TestModel> = .init()
    // TODO: Add documentation to location in firemodel schema.
    dynamic var location: Pring.GeoPoint?
    // TODO: Add documentation to colors in firemodel schema.
    dynamic var colors: [String]?
    // TODO: Add documentation to directions in firemodel schema.
    dynamic var directions: [TestEnum]?
    // TODO: Add documentation to models in firemodel schema.
    dynamic var models: [TestStruct]?
    // TODO: Add documentation to refs in firemodel schema.
    dynamic var refs: [Pring.AnyReference] = .init()
    // TODO: Add documentation to modelRefs in firemodel schema.
    dynamic var modelRefs: [Pring.Reference<TestTimestamps>] = .init()
    // TODO: Add documentation to meta in firemodel schema.
    dynamic var meta: [String: Any] = [:]
    // TODO: Add documentation to metaStrs in firemodel schema.
    dynamic var metaStrs: [String: String] = [:]
    // TODO: Add documentation to direction in firemodel schema.
    dynamic var direction: TestEnum?
    // TODO: Add documentation to testFile in firemodel schema.
    dynamic var testFile: Pring.File?
    // TODO: Add documentation to url in firemodel schema.
    dynamic var url: URL?
    // TODO: Add documentation to nested in firemodel schema.
    dynamic var nested: TestStruct?
    // TODO: Add documentation to nested_collection in firemodel schema.
    dynamic var nestedCollection: Pring.NestedCollection<TestModel> = []

    override func encode(_ key: String, value: Any?) -> Any? {
        switch key {
        case "direction":
            return self.direction?.firestoreValue
        case "nested":
            return self.nested?.rawValue
        default:
            break
        }
        return nil
    }

    override func decode(_ key: String, value: Any?) -> Bool {
        switch key {
        case "direction":
            self.direction = TestEnum(firestoreValue: value)
        case "nested":
          if let value = value as? [String: Any] {
            self.nested = TestStruct(id: self.id, value: value)
            return true
          }
        default:
            break
        }
        return false
    }
}

// TODO: Add documentation to TestTimestamps in firemodel schema.
@objcMembers class TestTimestamps: Pring.Object {
    static var testTimestampsId: String = ""
    override class var path: String { return "timestamps/\(testTimestampsId)" }
}

// TODO: Add documentation to Test in firemodel schema.
@objcMembers class Test: Pring.Object {
    
    // TODO: Add documentation to direction in firemodel schema.
    dynamic var direction: TestEnum?

    override func encode(_ key: String, value: Any?) -> Any? {
        switch key {
        case "direction":
            return self.direction?.firestoreValue
        default:
            break
        }
        return nil
    }

    override func decode(_ key: String, value: Any?) -> Bool {
        switch key {
        case "direction":
            self.direction = TestEnum(firestoreValue: value)
        default:
            break
        }
        return false
    }
}
