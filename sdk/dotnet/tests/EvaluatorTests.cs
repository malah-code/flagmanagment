using Xunit;
using FlagManagment.Sdk;

namespace FlagManagment.Sdk.Tests
{
    public class EvaluatorTests
    {
        [Fact]
        public void BucketUser_ReturnsValueBetweenZeroAndNinetyNine()
        {
            for (int i = 0; i < 1000; i++)
            {
                int bucket = Evaluator.BucketUser($"user-{i}");
                Assert.InRange(bucket, 0, 99);
            }
        }

        [Fact]
        public void BucketUser_IsDeterministicForSameKey()
        {
            int b1 = Evaluator.BucketUser("user_consistent_target");
            int b2 = Evaluator.BucketUser("user_consistent_target");
            Assert.Equal(b1, b2);
        }
    }
}
