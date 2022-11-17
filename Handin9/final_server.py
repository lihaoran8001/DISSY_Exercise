import requests
from http import cookiejar
from Crypto.Random import get_random_bytes
import string
import secrets
from Crypto.Util.Padding import pad,unpad
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

BLOCK_SIZE = 16

def xor(xs,ys):
    return bytes(x ^ y for x,y in zip(xs,ys))




## -----------------------------local-----------------------------

check_url='http://127.0.0.1:5000/check/'
url='http://127.0.0.1:5000/'

## -----------------------------website-----------------------------

# check_url='https://cbc-rsa.netsec22.dk:8000/check/'
# url='https://cbc-rsa.netsec22.dk:8000/'


##------------Part 1. use oracle padding to cover the secret
req=requests.get(url,verify=False)
cookies = requests.utils.dict_from_cookiejar(req.cookies)
#print(req.status_code)
#print(req.content)
#print(cookies)
token=cookies['authtoken']

ciphertext=bytes.fromhex(token)
print("ciphertext: ",ciphertext,"LEN: ",len(ciphertext))
blocks = [ciphertext[i:i+BLOCK_SIZE] for i in range(0, len(ciphertext), BLOCK_SIZE)]
ivv=blocks[0]
print("Now start oracle padding to get the plaintext")
print("Please wait with patience...")
result=b''
ivvv=b''
ivvv +=ivv
for ct in blocks[1:]:
    #start oracle padding
    zeroing_iv = [0] * BLOCK_SIZE
    for pad_val in range(1,1+BLOCK_SIZE):
        padding_iv=[pad_val ^ b for b in zeroing_iv]

        for candidate in range(256):
            padding_iv[-pad_val]=candidate
            iv = bytes(padding_iv)
            req_cookie=iv+ct
            #check status using the req_cookie
            cookies={}
            cookies['authtoken']=req_cookie.hex()
            req=requests.get(check_url,cookies=cookies,verify=False)
            # print(pad_val,candidate,len(req.content),req.content)
            if len(req.content)>50 or len(req.content)<20:
                # print("can:",candidate)
                break
        zeroing_iv[-pad_val] = candidate ^ pad_val
        # print("pad_val",pad_val)
    pt=bytes(iv_byte^dec_byte for iv_byte,dec_byte in zip(ivv,zeroing_iv))
##    print("zeroing_iv:",zeroing_iv)
    print("pt:",pt)
    # zeroing_iv from list to hex
    mid_iv=bytes(zeroing_iv)
    ivvv += mid_iv
    result += pt
    ivv=ct
    
print("covered secret:",result)