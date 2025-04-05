import csv

provincias = []
cantones = []
distritos = []


def capitalizador(inputString):
    lista = inputString.split(" ")
    outputString = ""
    for i in lista:
        outputString = outputString+i.capitalize()+" "
    outputString = outputString.rstrip()
    return (outputString)


with open("Nomenclator.csv", newline='') as distelec:
    lector = csv.reader(distelec)
    for row in lector:
        newProvincia = [row[0][:-4], capitalizador(row[3])]
        if newProvincia not in provincias:
            provincias.append(newProvincia)
        newCanton = [row[0][:-2], capitalizador(row[5])]
        if newCanton not in cantones:
            cantones.append(newCanton)
        newDistrito = [row[0], capitalizador(row[7])]
        if newDistrito not in distritos:
            distritos.append(newDistrito)

provincias.pop(0)
cantones.pop(0)
distritos.pop(0)

for i in provincias:
    print(f'INSERT INTO Provincias VALUES ({i[0]},"{i[1]}");')

for i in cantones:
    print(f'INSERT INTO Cantones VALUES ({i[0]},"{i[1]}",{i[0][:-2]});')

for i in distritos:
    print(f'INSERT INTO Distritos VALUES ({i[0]},"{i[1]}",{i[0][:-2]});')
