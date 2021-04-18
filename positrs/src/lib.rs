use softposit::P32;


#[no_mangle] pub extern "C" fn positadd(a: u32, b:u32) -> u32{
    let a = P32::from_bits(a);
    let b = P32::from_bits(b);
    let c = a+b;
    c.to_bits()
}

#[no_mangle] pub extern "C" fn positmul(a: u32, b:u32) -> u32{
    let a = P32::from_bits(a);
    let b = P32::from_bits(b);
    let c = a*b;
    c.to_bits()
}

#[no_mangle] pub extern "C" fn positdiv(a: u32, b:u32) -> u32{
    let a = P32::from_bits(a);
    let b = P32::from_bits(b);
    let c = a/b;
    c.to_bits()
}

#[no_mangle] pub extern "C" fn positsqrt(a: u32) -> u32{
    let a = P32::from_bits(a);
    let c = a.sqrt();
    c.to_bits()
}