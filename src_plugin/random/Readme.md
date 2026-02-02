# RANDOM

## PRÉSENTATION

Permet de simuler des succès ou des échecs d'authentification.  
Une option de suspension pour simuler une latence est disponible.

## OPTIONS

### **rand :**

Définir la durée en seconde de la suspension d'exécution.

    0 :
        Min
    
    120 :
        Max

    0 :
        Défaut

## EXEMPLE

Définir un temps d'arrêt de 5 secondes :

    [AUTH]
    mech = ["RANDOM"]

    [PLUGIN.RANDOM]
    rand = 5
