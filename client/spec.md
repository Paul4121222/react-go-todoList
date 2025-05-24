# ToDoList App 規格

## 1. 輸入區
- 一個文字輸入框，讓使用者輸入新的 todo 內容。
- 輸入框旁邊有一個「+」按鈕，點擊後會將輸入框的內容新增到 todoList，並清空輸入框。

## 2. Todo List 顯示區
- 以列表方式顯示所有 todo 項目。
- 每個 todo 項目需顯示：
  - 內容文字
  - 狀態（inprogress 或 done）
  - 狀態切換按鈕（例如：勾選框或按鈕，切換 inprogress/done）
  - 垃圾桶圖示按鈕，點擊可刪除該 todo

## 3. 狀態顯示
- 狀態為 inprogress 時，內容正常顯示。
- 狀態為 done 時，內容文字需加上刪除線（strike-through）。

## 4. 互動行為
- 新增：輸入內容並點擊「+」按鈕，新增 todo。
- 刪除：點擊垃圾桶圖示，刪除該 todo。
- 狀態切換：點擊狀態切換按鈕，todo 狀態在 inprogress/done 間切換。

## 5. UI/UX
- 輸入框與「+」按鈕同一行。
- Todo 列表整齊排列，狀態切換與刪除按鈕明顯易用。
- done 狀態的 todo 內容有刪除線，視覺上明顯區分。
