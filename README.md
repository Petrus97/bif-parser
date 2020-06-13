# Inferenza su reti bayesiane
Progetto per l'esame di intelligenza artificiale all'Università di Firenze.<br> Inferenza su reti bayesiane applicando l'algoritmo usato da Hugin Educational su Junction Tree.
## Esecuzione del codice
```
git clone
go get -d
cd example_test/
go test
```
# Riferimenti
Maggior parte del lavoro è stato svolto seguendo 
- [Probabilistic Graphical Models, Principles and Techniques By Daphne Koller and Nir Friedman](https://mitpress.mit.edu/books/probabilistic-graphical-models) per gli algoritmi utilizzati
- [Introduction to Bayesian Networks, F.V.Jensen](https://www.amazon.com/Introduction-Bayesian-Networks-Finn-Jensen/dp/0387915028) per gli esempi utilizzati
- [bnlear repository](https://www.bnlearn.com/bnrepository/) per le reti in formato BIF
- [ebay/bayesian-network](https://github.com/eBay/bayesian-belief-networks/blob/master/bayesian/examples/bif/bif_parser.py) per la realizzazione del parser
- [regex tutorial](https://github.com/StefanSchroeder/Golang-Regex-Tutorial) per imparare ad usare le regular expression con Go