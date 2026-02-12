# register

1. joyn.de
2. email
   1. mailsac.com fail
   2. gmail.com pass
   3. mail.tm?
3. password
   - Das Passwort muss mindestens 1 Klein- und Großbuchstaben sowie eine Zahl
   enthalten (The password must contain at least one lowercase and one uppercase
   letter, as well as one number)
4. Männlich (male)
5. day
6. month
7. year
8. save
9. PLUS
10. try now for free
11. enable JavaScript
12. first name
13. surname
14. Zahlungsmethode wählen (select payment method)
15. credit card
16. Karteninhaber (cardholder)
17. card number
18. CVC
19. month
20. year
21. I expressly agree to the execution of the contract
22. ordering the payment

## hexacard

https://hexacard.me/faq

## amazon pay

~~~json
POST https://p7.billwerk.com/api/v1/CustomerSelfService/Finalize HTTP/2.0

{
  "returnUrl": "https://www.joyn.de/abo/v2/willkommen?productVariantId=de-plus-web-monthly&interactivePayment=true&pactasTransactionId=69100f1a8d6ed9975022459e&secret=0sYLF0qPunj8F3s5FoVj5g&trigger=Payment&language=de-DE&amazonCheckoutSessionId=bfd7aff5-c6db-4f6c-805f-51135e87e5ba"
}

HTTP/2.0 200 

{
  "Error": {
    "Code": "InvalidCountry",
    "Message": "Payment failed!",
    "Details": "CheckoutSession cannot be completed: the billing address country US is not included in the allowed whitelist. Allowed countries: DE"
  }
}
~~~

## prompt

a service offers these payment methods

1. visa
2. mastercard
3. american express
4. paypal
5. amazon pay

however I get this when using Amazon Pay:

CheckoutSession cannot be completed: the billing address country US is not
included in the allowed whitelist. Allowed countries: DE

and similar errors with the other methods. how can I get one of these payment
methods? offer one option at a time and I will review the options one by one

## Option 1: Use a Virtual Address Service in Germany

the BIN for my cards will still be US so this is not a valid option

## Option 2: Use a Multi-Currency Account or Virtual Card Service

neither Wise nor Revolut offer German BIN cards to my understanding, if I am
wrong correct me

## Option 3: Use a specialized Virtual Credit Card (VCC) Provider

this is not a valid option as you have not given any example providers

## Option 4 (Revised): Obtain a European Prepaid Card from VIABUY

this is not a valid option as that URL is dead

## Option 5: Use the PayCenter SupremaCard Mastercard

it seems German ID is required, which would make this an invalid option

https://supremacard.de

## Option 6: Use a Service-Specific Gift Card or Voucher

you have ignored the original prompt which does not include this option as a
payment method

## Option 7 (Revised): Use an Assisted Purchase Service from a German Package Forwarder

https://mygermany.com
