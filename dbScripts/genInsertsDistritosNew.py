import csv

provincias = []
cantones = []
distritos = []

# capitaliza la primera letra de todas las palabras en un string


def capitalizador(inputString):
    lista = inputString.split(" ")
    outputString = ""
    for i in lista:
        outputString = outputString+i.capitalize()+" "
    outputString = outputString.rstrip()
    return (outputString)


# abrimos el csv que contiene los datos
with open("Nomenclator.csv", newline='') as distelec:
    lector = csv.reader(distelec)
    # empezamos a guardar los datos de todas las filas en las listas
    for row in lector:
        # añade nueva provincia si ya no ha sido añadida previamente
        newProvincia = [row[0][:-4], capitalizador(row[3])]
        if newProvincia not in provincias:
            provincias.append(newProvincia)
        # añade nuevo canton si ya no ha sido añadido previamente
        newCanton = [row[0][:-2], capitalizador(row[5])]
        if newCanton not in cantones:
            cantones.append(newCanton)
        # añade nuevo distrito
        newDistrito = [row[0], capitalizador(row[7])]
        if newDistrito not in distritos:
            distritos.append(newDistrito)

# eliminamos la primera fila que contiene basura
provincias.pop(0)
cantones.pop(0)
distritos.pop(0)

# generamos los inserts
for i in provincias:
    print(f'INSERT INTO Provincias VALUES ({i[0]},"{i[1]}");')

for i in cantones:
    print(f'INSERT INTO Cantones VALUES ({i[0]},"{i[1]}",{i[0][:-2]});')

for i in distritos:
    print(f'INSERT INTO Distritos VALUES ({i[0]},"{i[1]}",{i[0][:-2]});')
