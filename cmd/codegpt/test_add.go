package main

func testAddedCodeReviewCalled() (int, error) {
	return 1, nil
}

func testAddedCodeReview() {
	a, err := testAddedCodeReviewCalled()
	_ = err
	print(a)
}
