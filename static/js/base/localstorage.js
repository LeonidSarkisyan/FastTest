const TESTS_KEY = "tests"

export function AddTestID(testID, index) {
    let tests = JSON.parse(localStorage.getItem(TESTS_KEY))

    if (tests) {
        if (InTests(testID)) {

        } else {
            tests.push({
                testID,
                index
            })
            localStorage.setItem(TESTS_KEY, JSON.stringify(tests))
        }
    } else {
        localStorage.setItem(TESTS_KEY, JSON.stringify([{
            testID,
            index
        }]))
    }
}

function InTests(testID) {
    let tests = JSON.parse(localStorage.getItem(TESTS_KEY))

    for (let test of tests) {
        if (test.testID === testID) {
            return true
        }
    }

    return false
}

export function ChangeIndex(testID, newIndex) {
    let tests = JSON.parse(localStorage.getItem(TESTS_KEY))

    let index = 0;
    for (let test of tests) {
        if (test.testID === testID) {
            console.log(test)
            tests[index].index = newIndex
        }
        index++
    }

    localStorage.setItem(TESTS_KEY, JSON.stringify(tests))
}

export function GetIndex(testID) {
    let tests = JSON.parse(localStorage.getItem(TESTS_KEY))

    let index = 0;
    for (let test of tests) {
        if (test.testID === testID) {
            return test.index
        }
        index++
    }

    return 0
}