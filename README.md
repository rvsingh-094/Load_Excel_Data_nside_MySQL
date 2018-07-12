# LOAD MULTIPLE EXCEL FILE INSIDE MYSQL TABLE

This ExcelFileUploader program is used to upload multiple file at a time inside the mysql table.
We just need to pass the all column name and folder path.

By using help of go routine and channel , we use dump all data inside the db.

Benefit :
		1. If we need to perform any operation on particular row or file in that case using go_routine and channel we can solve this type of problem.<br>
		2. We can upload larger number of file using different database at a given time by just distribute the work inside the go_routine.<br>
		3. So by using this scenario we can also understand easily how the go_routine and concurrency program work.


Use Case : 
		Step 1 : Add this program inside the project directory.<br>
		Step 2 : Call this TraverseAllDirectoryForExcel(<DbConnection>, <filepath>)<br>
		Step 3 : Pass the table name inside the TraverseAllDirectoryForExcel and set the query (right now i did it manually but next version just pass the table name if its exist then build query other create table based on excel first header,)<br>
		Step 5: If customization required in that case do customization	before inserting the row inside the channel or insert function (two level customization available)<br>
		Step 6: if excel file level customization then control it while traversing the file and based on that call customize go_routine.<br>
		Step 4 : Its automatically create go_routine , channel and insert all data.
		
			
