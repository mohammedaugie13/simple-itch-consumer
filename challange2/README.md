# Challange 2

For this refactor task i did several changes:
1. I put the summary map `summary := make(map[string]*Summary)`  outside ohlc function. So summary map will not as product of return funcion but instead we pass it as reference map so the OHLC function will update that.
2. GetListedIndex(stockCode string, indexes []IndexMember). Because we dont need append the index in main looping. And we just need one time to updaate the index for every stock code.
3. UpdateRecord(summary *Summary, incomingData ChangeRecord) is the main logic i refactor it using switch case instead if else.