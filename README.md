# auction
An auction system

Constraints:
- Information is not persistant.
- Auctioner needs to be up & running before bidders.
- 1 Auctioner & 3 Bidders. (Bidders can be increased easily)

Steps to run:
- greedy/auctioner 
   go run main.go 
- greedy/bidder
   go run main.go
- greedy/bidder1        	
   go run main.go
- greedy/bidder2        	
   go run main.go   

Endpoints:
- Fetch all registered bidder endpoints
GET : http://localhost:8080/bidder

- Run an auction
POST: http://localhost:8080/auction

