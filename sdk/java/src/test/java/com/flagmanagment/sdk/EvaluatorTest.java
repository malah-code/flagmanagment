package com.flagmanagment.sdk;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

public class EvaluatorTest {

    @Test
    public void testBucketUserDeterministicRange() {
        for (int i = 0; i < 1000; i++) {
            int bucket = Evaluator.bucketUser("user-" + i);
            assertTrue(bucket >= 0 && bucket < 100, "Bucket must be between 0 and 99");
        }
    }

    @Test
    public void testBucketUserConsistent() {
        int b1 = Evaluator.bucketUser("user_consistent_target");
        int b2 = Evaluator.bucketUser("user_consistent_target");
        assertEquals(b1, b2, "Same user key must always yield the identical bucket");
    }
}
