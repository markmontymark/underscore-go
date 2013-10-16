use strict;
use warnings;

check($_) for @ARGV;
exit;

sub check 
{
	my $f = shift;

	my %funcs = ();
	my %oops = ();
	open F,$f or die "Can't read file $f, $!\n";
	while( <F> )
	{
		chomp ;
		if( /^\s*func\s+([^\s(]+)\(/ ) {
			#print "func $1\n";
			$funcs{$1} = $_;
		}
		elsif( /^\s*func\s+\(\s*this\s+\*\s*Underscore\s*\)\s*([^\s(]+)\(/ ) {
			#print "oop  $1\n";
			$oops{$1} = $_;
		}
	}
	close F;

	for (sort keys %funcs){
		unless( exists $oops{$_} ){
			print "func (this \*Underscore) $_ () *Underscore {
// $funcs{$_}
	if this.ischained {
		return New( $_\( this.wrapped ))
	}
	this.wrapped = $_\( this.wrapped )
	return this
}

";
		}
	}
	for (sort keys %oops){
		unless( exists $funcs{$_} ){
			print "missing func for $_\n";
		}
	}
}

