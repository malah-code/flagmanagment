import XCTest
@testable import FlagManagment

final class EvaluatorTests: XCTestCase {
    func testBucketUserRange() throws {
        for i in 0..<1000 {
            let bucket = Evaluator.bucketUser("user-\(i)")
            XCTAssertGreaterThanOrEqual(bucket, 0)
            XCTAssertLessThan(bucket, 100)
        }
    }

    func testBucketUserDeterministic() throws {
        let b1 = Evaluator.bucketUser("user_consistent_target")
        let b2 = Evaluator.bucketUser("user_consistent_target")
        XCTAssertEqual(b1, b2)
    }
}
